package handlers

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"html/template"
	"net/http"
)

func AllMetricsHTML(getter repository.CollectionGetter) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./internal/server/templates/index.tmpl"))
		rw.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(rw, getter.GetCollection())
		if err != nil {
			String(rw, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
