package errata

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/flosch/pongo2/v5"
	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	//go:embed web/*
	web embed.FS
)

type Server struct {
	File    string
	Package string

	db DataSource
}

func NewServer(path, pkg string) (*Server, error) {
	source, err := NewHCLDatasource(path)
	if err != nil {
		return nil, err
	}

	return &Server{
		File:    path,
		Package: pkg,
		db:      source,
	}, nil
}

func (s *Server) List(w http.ResponseWriter, req *http.Request) {
	data := pongo2.Context{
		"Package":       s.Package,
		"Errors":        s.db.List(),
		"Options":       s.db.Options(),
		"LastUpdatedAt": time.Now().Format(time.RFC3339),
	}
	renderMarkdown(s.db)

	path := "web/list.gohtml"
	_, err := web.Open(path)
	if err != nil {
		s.errorHandler(NewTemplateNotFoundErr(err), w)
		return
	}

	b, err := web.ReadFile(path)
	if err != nil {
		s.errorHandler(NewTemplateNotReadableErr(err), w)
	}

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	pongo2.RegisterFilter("markdown", filterMarkdown(md))

	tmpl, err := pongo2.FromBytes(b)
	if err != nil {
		s.errorHandler(NewTemplateSyntaxErr(err), w)
		return
	}

	if err := tmpl.ExecuteWriter(data, w); err != nil {
		s.errorHandler(NewTemplateExecutionErr(err), w)
		return
	}
}

func (s *Server) Item(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	code := params["code"]

	erratum, ok := s.db.FindByCode(code)
	if !ok {
		s.errorHandler(NewServeUnknownCodeErr(nil, code), w)
		return
	}

	data := pongo2.Context{
		"Package":       s.Package,
		"Error":         erratum,
		"Code":          code,
		"Options":       s.db.Options(),
		"LastUpdatedAt": time.Now().Format(time.RFC3339),
	}
	renderMarkdown(s.db)

	path := "web/single.gohtml"
	_, err := web.Open(path)
	if err != nil {
		s.errorHandler(NewTemplateNotFoundErr(err), w)
		return
	}

	b, err := web.ReadFile(path)
	if err != nil {
		s.errorHandler(NewTemplateNotReadableErr(err), w)
	}

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	pongo2.RegisterFilter("markdown", filterMarkdown(md))

	tmpl, err := pongo2.FromBytes(b)
	if err != nil {
		s.errorHandler(NewTemplateSyntaxErr(err), w)
		return
	}

	if err := tmpl.ExecuteWriter(data, w); err != nil {
		s.errorHandler(NewTemplateExecutionErr(err), w)
		return
	}
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

func (s *Server) errorHandler(err error, w http.ResponseWriter) {
	// TODO: pick HTTP status code if available
	// 		 maybe with the pattern errata.HTTPStatusExtractor(err, default=http.StatusInternalServerError)?
	http.Error(w, err.Error(), http.StatusInternalServerError)
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
	r.HandleFunc("/code/{code}", srv.Item)

	return http.ListenAndServe("localhost:33707", r)
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
