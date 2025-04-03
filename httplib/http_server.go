package httplib

import (
	"context"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const readHeaderTimeout = 3 * time.Second

type Server interface {
	Start()
	Stop()
}

type api struct {
	server          *http.Server
	port            string
	shutdownTimeout time.Duration
	done            chan struct{}
}

func NewHTTPServer(router http.Handler, port string, shutdownTimeout time.Duration) Server {
	server := &http.Server{
		Addr:              port,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &api{
		server:          server,
		port:            port,
		shutdownTimeout: shutdownTimeout,
		done:            make(chan struct{}),
	}
}

func (a *api) Start() {
	go func() {
		defer close(a.done)

		log.WithField("address", a.port).Info("starting http server")

		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Fatal("http listen and serve")
		}

		log.Info("http server stopped listening")
	}()
}

func (a *api) Stop() {
	log.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	if err := a.server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("http shutdown error")
	}

	log.WithField("port", a.port).Info("http server stopped")

	<-a.done
	cancel()
}
