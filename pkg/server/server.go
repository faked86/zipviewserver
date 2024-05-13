package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Ext: {{.Ext}}</title>
	</head>
	<body>
		<h1>Ext: {{.Ext}}</h1>
		{{if .Items}}
			<ul>
				{{range .Items}}
					<li>{{ . }}</li>
				{{end}}
			</ul>
		{{else}}
			<p>No such files</p>
		{{end}}
	</body>
</html>`

func NewServer(port int, ctx context.Context, ext string, fileNames []string) *http.Server {
	r := chi.NewRouter()
	r.Get("/", indexHandler(ext, fileNames))
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	return httpServer
}

func indexHandler(ext string, fileNames []string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("New request")
		t, err := template.New("page").Parse(tpl)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		data := struct {
			Ext   string
			Items []string
		}{
			Ext:   ext,
			Items: fileNames,
		}
		t.Execute(w, data)
	}
}
