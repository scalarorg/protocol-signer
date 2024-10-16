package signerservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/protocol-signer/config"
	"github.com/scalarorg/protocol-signer/internals/signer/handlers"
	"github.com/scalarorg/protocol-signer/internals/signer/middlewares"
	"github.com/scalarorg/protocol-signer/packages/btc"
	"github.com/scalarorg/protocol-signer/packages/evm"
)

type SigningServer struct {
	httpServer *http.Server
	handler    *handlers.Handler
}

func (a *SigningServer) SetupRoutes(r *chi.Mux) {
	handler := a.handler
	r.Post("/v1/sign-unbonding-tx", registerHandler(handler.SignUnbonding))
}

func New(ctx context.Context, cfg *config.ParsedConfig, btcSigner *btc.PsbtSigner) (*SigningServer, error) {
	r := chi.NewRouter()

	// TODO: Add middlewares
	// r.Use(middlewares.CorsMiddleware(cfg))
	r.Use(middlewares.TracingMiddleware)
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.ContentLengthMiddleware(int64(cfg.ServerConfig.MaxContentLength)))
	// TODO: TLS configuration if server is to be exposed over the internet, if it supposed to
	// be behind some reverse proxy like nginx or cloudflare, then it's not needed.
	// Probably it needs to be configurable

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port),
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
		ReadTimeout:  cfg.ServerConfig.ReadTimeout,
		IdleTimeout:  cfg.ServerConfig.IdleTimeout,
		Handler:      r,
	}
	// Create evm clients
	evmClients := make([]evm.EvmClient, len(cfg.EvmConfigs))
	for i, evmConfig := range cfg.EvmConfigs {
		evmClients[i] = *evm.NewEvmClient(evmConfig)
	}
	h, err := handlers.NewHandler(evmClients, btcSigner)
	if err != nil {
		log.Fatal().Err(err).Msg("error while setting up handlers")
	}

	server := &SigningServer{
		httpServer: srv,
		handler:    h,
	}
	server.SetupRoutes(r)
	return server, nil
}

func (s *SigningServer) Start() error {
	log.Info().Msgf("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *SigningServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
