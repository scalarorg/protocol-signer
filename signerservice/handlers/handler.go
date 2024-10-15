package handlers

import (
	"context"
	"fmt"
	"net/http"

	m "github.com/scalarorg/protocol-signer/observability/metrics"
	s "github.com/scalarorg/protocol-signer/signerapp"
)

type Handler struct {
	t    string
	evms []ExternalEvmClient
	s    *s.SignerApp
	m    *m.CovenantSignerMetrics
}

type Result struct {
	Data   interface{}
	Status int
}

type PublicResponse[T any] struct {
	Data T `json:"data"`
}

func NewResult[T any](data T) *Result {
	res := &PublicResponse[T]{Data: data}
	return &Result{Data: res, Status: http.StatusOK}
}

func NewHandler(
	_ context.Context, t string, evms []ExternalEvmClient, s *s.SignerApp, m *m.CovenantSignerMetrics,
) (*Handler, error) {
	return &Handler{
		t:    t,
		evms: evms,
		s:    s,
		m:    m,
	}, nil
}

func (h *Handler) getEvmClient(chainName string) (ExternalEvmClient, error) {
	for _, evm := range h.evms {
		if evm.ChainName() == chainName {
			return evm, nil
		}
	}
	return nil, fmt.Errorf("evm client not found for chain name: %s", chainName)
}
