package server

import (
	"html/template"
	"net/http"
)

func (s *server) AllMetricsHTML(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./internal/server/templates/index.tmpl"))
	rw.Header().Set("Content-Type", "text/html")
	collection, err := s.storage.GetCollection()
	if err != nil {
		String(rw, http.StatusInternalServerError, err.Error())
		return
	}
	err = tmpl.Execute(rw, collection)
	if err != nil {
		String(rw, http.StatusInternalServerError, err.Error())
		return
	}
}
