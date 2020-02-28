package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

// Server containts the configured app server
type Server struct {
	Store  store.Store
	Server *http.Server
	Router *chi.Mux
	// Log *log.Logger
	// RateLimiter *RateLimiter
	// jobs
	// other cfg
}

// NewServer creates the new server
func NewServer(st store.Store) (*Server, error) {
	r := chi.NewRouter()

	s := &Server{
		Router: r,
		Store:  st,
	}

	// s.Log = log.NewLogger()

	return s, nil
}

// Start runs the HTTP server
func (s *Server) Start() (err error) {
	// TODO: read from cfg
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Content-Length", "Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma", "Authorization"},
		MaxAge:           86400,
		AllowCredentials: true,
		Debug:            false,
	})

	handler := corsWrapper.Handler(s.Router)

	// if ratelimit set it here...

	listenAddr := ":3001"

	s.Server = &http.Server{
		// ErrorLog: logger
		Handler:           handler,
		Addr:              listenAddr,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      7 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// autocert manager later...
	// if UseLetsEncrypt...

	go func() {
		if err = s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", listenAddr, err)
		}
	}()
	log.Printf("server started")

	gracefullShutdown(s.Server)

	return
}

func gracefullShutdown(srv *http.Server) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %+s", err)
	}
	log.Fatalf("server is shutting down")
}
