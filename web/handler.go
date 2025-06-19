package web

import (
	"net/http"
	"text/template"

	goreddit "github.com/benitopedro13/go-reddit"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	*chi.Mux

	store goreddit.Store
}

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Mux.Use(middleware.Logger)

	h.Mux.Route("/threads", func(r chi.Router) {
		r.Get("/", h.ThreadsList())
	})

	return h
}

const threadsListHTML = `
	<h1>Threads</h1>
	<dl>
	{{range .Threads}}
		<dt><strong>{{.Title}}</strong></dt>
		<dd>{{.Description}}</dd>
	{{end}}
	</dl>
	`

func (h *Handler) ThreadsList() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}

	tmpl := template.Must(template.New("").Parse(threadsListHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Threads: tt})
	}
}
