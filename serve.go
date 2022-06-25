package errata

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/flosch/pongo2/v5"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	//go:embed web/*
	web embed.FS
)

type Server struct {
	File string

	logger   log.Logger
	db       DataSource
	idx      bleve.Index
	bindAddr string
}

func NewServer(logger log.Logger, cfg WebUIConfig) (*Server, error) {
	source, err := NewHCLDatasource(cfg.Source)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		File:     cfg.Source,
		db:       source,
		logger:   log.With(logger, "component", "web"),
		bindAddr: cfg.BindAddr,
	}

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	pongo2.RegisterFilter("markdown", filterMarkdown(md))

	err = srv.buildIndex()
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Server) buildIndex() error {
	idx, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return err
	}

	for code, e := range s.db.List() {
		if err := idx.Index(code, e); err != nil {
			return NewServeSearchIndexErr(err)
		}
	}

	s.idx = idx
	return nil
}

func (s *Server) Search(w http.ResponseWriter, req *http.Request) {
	term := req.FormValue("term")
	if strings.TrimSpace(term) == "" {
		s.errorHandler(w, NewServeSearchMissingTermErr(nil))
		return
	}

	query := bleve.NewMatchPhraseQuery(term)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := s.idx.Search(searchRequest)

	if len(searchResult.Hits) == 0 {
		s.render(w, "web/search-miss.gohtml", pongo2.Context{
			"Term": term,
		})
		return
	}

	errs := make(map[string]errorDefinition, len(searchResult.Hits))
	for _, e := range searchResult.Hits {
		err, _ := s.db.FindByCode(e.ID)
		errs[e.ID] = err
	}

	data := pongo2.Context{
		"Errors":        errs,
		"Options":       s.db.Options(),
		"LastUpdatedAt": time.Now().Format(time.RFC3339),
	}
	s.render(w, "web/list.gohtml", data)
}

func (s *Server) List(w http.ResponseWriter, _ *http.Request) {
	data := pongo2.Context{
		"Errors":        s.db.List(),
		"Options":       s.db.Options(),
		"LastUpdatedAt": time.Now().Format(time.RFC3339),
	}
	s.render(w, "web/list.gohtml", data)
}

func (s *Server) render(w http.ResponseWriter, path string, data pongo2.Context) {
	renderMarkdown(s.db)

	_, err := web.Open(path)
	if err != nil {
		s.errorHandler(w, NewTemplateNotFoundErr(err))
		return
	}

	b, err := web.ReadFile(path)
	if err != nil {
		s.errorHandler(w, NewTemplateNotReadableErr(err))
	}

	tmpl, err := pongo2.FromBytes(b)
	if err != nil {
		s.errorHandler(w, NewTemplateSyntaxErr(err))
		return
	}

	if err := tmpl.ExecuteWriter(data, w); err != nil {
		s.errorHandler(w, NewTemplateExecutionErr(err))
		return
	}
}

func (s *Server) Item(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	code := params["code"]

	erratum, ok := s.db.FindByCode(code)
	if !ok {
		s.errorHandler(w, NewServeUnknownCodeErr(nil, code))
		return
	}

	data := pongo2.Context{
		"Error":         erratum,
		"Code":          code,
		"Options":       s.db.Options(),
		"LastUpdatedAt": time.Now().Format(time.RFC3339),
	}

	s.render(w, "web/single.gohtml", data)
}

func filterMarkdown(md goldmark.Markdown) func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	return func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		var buf bytes.Buffer
		mdErr := md.Convert([]byte(in.String()), &buf)
		if mdErr != nil {
			return nil, &pongo2.Error{
				OrigError: NewMarkdownRenderErr(mdErr),
				Sender:    "filterMarkdown",
			}
		}

		return pongo2.AsSafeValue(buf.String()), nil
	}
}

func (s *Server) errorHandler(w http.ResponseWriter, err error) {
	LogError(s.logger, err)

	statusCode := http.StatusInternalServerError
	if e, ok := err.(HTTPStatusCodeExtractor); ok {
		if c, cerr := strconv.ParseInt(e.GetHttpStatusCode(), 10, 16); cerr == nil {
			statusCode = int(c)
		}
	}

	http.Error(w, err.Error(), statusCode)
}

type HTTPStatusCodeExtractor interface {
	GetHttpStatusCode() string
}

func Serve(srv *Server) error {
	r := mux.NewRouter()

	webFS, err := fs.Sub(web, "web")
	if err != nil {
		// TODO wrap error
		return err
	}
	r.PathPrefix("/assets/").Handler(http.FileServer(http.FS(webFS)))
	r.HandleFunc("/favicon.ico", http.FileServer(http.FS(webFS)).ServeHTTP)
	r.HandleFunc("/", srv.List)
	r.HandleFunc("/code/{code}", srv.Item).Methods(http.MethodGet)
	r.HandleFunc("/search", srv.Search).Methods(http.MethodGet)

	level.Info(srv.logger).Log("msg", "web UI started", "bind-addr", srv.bindAddr)
	return http.ListenAndServe(srv.bindAddr, r)
}

func renderMarkdown(source DataSource) {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	for _, e := range source.List() {
		if err := md.Convert([]byte(fmt.Sprintf("%s", e.Guide)), &buf); err != nil {
			fmt.Println(NewMarkdownRenderErr(err))
		}

		e.Guide = buf.String()
	}
}
