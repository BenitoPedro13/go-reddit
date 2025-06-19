package web

import (
	"net/http"
	"text/template"

	goreddit "github.com/benitopedro13/go-reddit"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
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
		r.Get("/new", h.ThreadsCreate())
		r.Post("/", h.ThreadsStore())
		r.Post("/{id}/delete", h.ThreadsDelete())
	})

	return h
}

const threadsListHTML = `
	<h1>Threads</h1>
	<dl>
	{{range .Threads}}
		<dt><strong>{{.Title}}</strong></dt>
		<dd>{{.Description}}</dd>
		<dd>
			<form method="post" action="/threads/{{.ID}}/delete">
				<button type="submit">Delete</button>
			</form>
		</dd>
	{{end}}
	</dl>
	<a href="/threads/new">Create thread</a>
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

const threadsCreateHTML = `
	<h1>Create Thread</h1>
	<form method="post" action="/threads">
		<table>
			<tr>
				<td><label for="title">Title</label></td>
				<td><input type="text" name="title" id="title" /></td>
			</tr>
			<tr>
				<td><label for="description">Description</label></td>
				<td><input type="text" name="description" id="description" /></td>
			</tr>
		</table>
		<button type="submit">Create thread</button>
	</form>
	`

func (h *Handler) ThreadsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(threadsCreateHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) ThreadsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		err := h.store.CreateThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusSeeOther)
	}
}

func (h *Handler) ThreadsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.store.DeleteThread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/threads", http.StatusSeeOther)
	}
}
