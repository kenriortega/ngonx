package cli

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kenriortega/ngonx/pkg/logger"
)

type server struct {
	*http.Server
}

func NewServer(host string, port int) *server {

	s := &http.Server{
		Handler: nil,
		Addr:    fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	return &server{s}
}

func NewServerSSL(host string, port int) *server {
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	s := &http.Server{
		Handler: nil,
		Addr:    fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	return &server{s}
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *server) Start() {
	logger.LogInfo("starting server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError(
				fmt.Sprintf("could not listen on %s due to %s", srv.Addr, err.Error()),
			)
		}
	}()
	logger.LogInfo(fmt.Sprintf("server is ready to handle requests %s", srv.Addr))
	srv.gracefulShutdown()
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *server) StartSSL(crt, key string) {
	logger.LogInfo("starting server...")

	go func() {
		if err := srv.ListenAndServeTLS(crt, key); err != nil && err != http.ErrServerClosed {
			logger.LogError(
				fmt.Sprintf("could not listen on %s due to %s", srv.Addr, err.Error()),
			)
		}
	}()
	logger.LogInfo(fmt.Sprintf("server is ready to handle requests %s", srv.Addr))
	srv.gracefulShutdown()
}

func (srv *server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	logger.LogInfo(fmt.Sprintf("server is shutting down %s", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		logger.LogError(
			fmt.Sprintf("could not gracefully shutdown the server %s", err.Error()),
		)
	}
	logger.LogInfo("server stopped")
}
