package errata

import (
	"embed"
	"html/template"
	"net/http"
)

var (
	//go:embed web/*
	web embed.FS
)

type Server struct {
	File    string
	Package string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	source, err := NewHCLDatasource(s.File)
	if err != nil {
		s.errorHandler(err, w)
		return
	}

	renderMarkdown(source)

	tmplData := Tmpl{
		Package: s.Package,
		Errors:  source.List(),
	}

	tmpl, err := template.New("index.gohtml").
		ParseFS(web, "web/*")

	if err != nil {
		s.errorHandler(NewTemplateSyntaxErr(err), w)
		return
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		s.errorHandler(NewTemplateExecutionErr(err), w)
		return
	}
}

func (s *Server) errorHandler(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func Serve(srv *Server) error {
	return http.ListenAndServe("localhost:8080", srv)
}

func renderMarkdown(source DataSource) {
	return
	// render markdown
	//var buf bytes.Buffer
	//md := goldmark.New(
	//	goldmark.WithRendererOptions(
	//		html.WithHardWraps(),
	//	),
	//)
	//
	//for _, e := range source.List() {
	//	if sol, ok := e.Definition["solution"]; ok {
	//		if err := md.Convert([]byte(fmt.Sprintf("%s", sol)), &buf); err != nil {
	//			fmt.Println(NewMarkdownRenderErr(err))
	//		}
	//
	//		e.Definition["solution"] = template.HTML(buf.String())
	//	}
	//}
}
