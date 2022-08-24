package server

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type server struct {
	address string
	router  *chi.Mux
	storage repository.Storager
	key     string
}

func NewServer(address string, storage repository.Storager, key string) *server {
	s := &server{
		address: address,
		storage: storage,
		key:     key,
	}

	s.router = chi.NewRouter()
	s.router.Use(middleware.StripSlashes)
	s.router.Use(middlewares.PackGzip)
	s.router.Use(middlewares.UnpackGzip)
	s.registerRoutes()

	return s
}

func (s *server) registerRoutes() {
	s.router.NotFound(NotFound)
	s.router.MethodNotAllowed(MethodNotAllowed)

	s.router.Get("/", s.AllMetricsHTML)
	s.router.Get("/value/{type}/{name}", s.GetMetricHandler)
	s.router.Post("/update/{type}/{name}/{value}", s.UpdateMetricHandler)
	s.router.Post("/value", s.GetMetricJSONHandler)
	s.router.Post("/update", s.UpdateMetricJSONHandler)
}

func (s *server) Run() {
	fmt.Printf("The server has just started on address [%s]\n", s.address)
	log.Fatal(http.ListenAndServe(s.address, s.router))
}
