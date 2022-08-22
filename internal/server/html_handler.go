package server

import (
	"html/template"
	"net/http"
)

func (s *server) AllMetricsHTML(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./internal/server/templates/index.tmpl"))
	rw.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(rw, s.storage.GetCollection())
	if err != nil {
		String(rw, http.StatusInternalServerError, err.Error())
		return
	}
}
