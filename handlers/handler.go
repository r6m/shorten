package handlers

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/r6m/shorten/store"
	"github.com/sirupsen/logrus"
)

type Handler func(w http.ResponseWriter, r *http.Request) any

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := h(w, r)
	if v == nil {
		return
	}

	switch v.(type) {
	case error:
		handleError(w, r, v.(error))
	default:
		render.Respond(w, r, v)
	}
}

func handler(h Handler) http.HandlerFunc {
	return h.ServeHTTP
}

// API is our handler wrapper
type API struct {
	r    chi.Router
	repo store.Store
}

// NewServer sets up api server routes
func NewServer(repo store.Store) *API {
	api := &API{
		repo: repo,
	}

	r := chi.NewRouter()
	api.r = r

	r.Post("/shorten", handler(api.shortenHandler))
	r.Get("/{key}", handler(api.redirectHandler))
	r.Get("/{key}/info", handler(api.infoHandler))

	return api
}

// ListenAndServe starts to listen on given addr
func (api *API) ListenAndServe(addr string) {
	log := logrus.WithField("component", "api")

	server := &http.Server{
		Addr:    addr,
		Handler: api.r,
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-quit:
			log.Infof("signal %s received. shutdown server", sig)
		case <-done:
			log.Infof("shutting down...")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.Shutdown(ctx)
	}()

	log.Printf("starting server on http://localhost%s", addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("http server failed to start")
	}
}
