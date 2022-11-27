package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/middlewares"
)

type Server struct {
	address string
	router  *chi.Mux
	storage storage.Storager
	key     string
}

func NewServer(address string, storage storage.Storager, key string) *Server {
	s := &Server{
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

func (s *Server) registerRoutes() {
	s.router.NotFound(NotFound)
	s.router.MethodNotAllowed(MethodNotAllowed)

	s.router.Get("/debug/pprof", pprof.Index)
	s.router.Get("/debug/pprof/profile", pprof.Profile)
	s.router.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	s.router.Get("/", s.AllMetricsHTML)
	s.router.Get("/ping", s.PingDBHandler)
	s.router.Get("/value/{type}/{name}", s.GetMetricHandler)
	s.router.Post("/update/{type}/{name}/{value}", s.UpdateMetricHandler)
	s.router.Post("/value", s.GetMetricJSONHandler)
	s.router.Post("/update", s.UpdateMetricJSONHandler)
	s.router.Post("/updates", s.UpdateMetricsBatchJSONHandler)
}

func (s *Server) Run() {
	log.Info().Str("address", s.address).Msg("The server has just started")

	if err := http.ListenAndServe(s.address, s.router); err != nil {
		log.Fatal().Err(err).Msg("Server error occurred")
	}
}
