package server

import (
	"context"
	"crypto/rsa"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/middlewares"
)

type HTTPServer struct {
	core    *http.Server
	router  *chi.Mux
	storage storage.Storager
	key     string
}

func newHTTPServer(address string, storage storage.Storager, key string, privateKey *rsa.PrivateKey, trustedSubnet string) *HTTPServer {
	s := &HTTPServer{
		storage: storage,
		key:     key,
	}

	s.router = chi.NewRouter()
	s.router.Use(middleware.StripSlashes)
	s.router.Use(middlewares.TrustedSubnet(trustedSubnet))
	s.router.Use(middlewares.PackGzip)
	s.router.Use(middlewares.UnpackGzip)
	s.router.Use(middlewares.Decrypt(privateKey))
	s.registerRoutes()

	s.core = &http.Server{
		Addr:    address,
		Handler: s.router,
	}

	return s
}

func (s *HTTPServer) registerRoutes() {
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

func (s *HTTPServer) Run() error {
	log.Info().Str("address", s.core.Addr).Msg("The server has just started")

	return s.core.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.core.Shutdown(ctx)
}
