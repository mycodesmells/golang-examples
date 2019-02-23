package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/satori/go.uuid"
)

func main() {
	svc := service{
		db: map[string]*quote{},
	}

	mux := chi.NewRouter()
	mux.Post("/quotes", svc.AddQuote)
	mux.Get("/quotes", svc.ListQuotes)
	mux.Get("/quotes/{id}", svc.GetQuote)
	mux.Delete("/quotes/{id}", svc.DeleteQuote)

	http.ListenAndServe(":8888", mux)
}

type quote struct {
	ID     string `json:"id,omitempty"`
	Quote  string `json:"quote,omitempty"`
	Author string `json:"author,omitempty"`
}

func (q *quote) Bind(req *http.Request) error {
	return nil
}

type service struct {
	db map[string]*quote
}

func (s *service) AddQuote(rw http.ResponseWriter, req *http.Request) {
	q := &quote{}
	if err := render.Bind(req, q); err != nil {
		render.PlainText(rw, req, err.Error())
		return
	}
	q.ID = uuid.NewV1().String()

	s.db[q.ID] = q

	render.JSON(rw, req, q)
}

func (s *service) ListQuotes(rw http.ResponseWriter, req *http.Request) {
	qq := make([]*quote, 0, len(s.db))
	for _, q := range s.db {
		qq = append(qq, q)
	}

	render.JSON(rw, req, qq)
}

func (s *service) GetQuote(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	q, ok := s.db[id]
	if !ok {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	render.JSON(rw, req, q)
}

func (s *service) DeleteQuote(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	if _, ok := s.db[id]; !ok {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}
	delete(s.db, id)

	render.NoContent(rw, req)
}
