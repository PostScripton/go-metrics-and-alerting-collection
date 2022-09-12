package server

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"net/http"
)

type server struct {
	address string
	router  *chi.Mux
	storage repository.Storager
	key     string
	pool    *pgxpool.Pool
}

func NewServer(address string, storage repository.Storager, key string, pool *pgxpool.Pool) *server {
	s := &server{
		address: address,
		storage: storage,
		key:     key,
		pool:    pool,
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
	s.router.Get("/ping", s.PingDBHandler)
	s.router.Get("/value/{type}/{name}", s.GetMetricHandler)
	s.router.Post("/update/{type}/{name}/{value}", s.UpdateMetricHandler)
	s.router.Post("/value", s.GetMetricJSONHandler)
	s.router.Post("/update", s.UpdateMetricJSONHandler)
	s.router.Post("/updates", s.UpdateMetricsBatchJSONHandler)
}

func (s *server) Run() {
	log.Info().Str("address", s.address).Msg("The server has just started")

	if err := http.ListenAndServe(s.address, s.router); err != nil {
		log.Fatal().Err(err).Msg("Server error occurred")
	}
}
