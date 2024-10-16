package handlers

import (
	"fmt"
	"net/http"

	"github.com/scalarorg/protocol-signer/packages/btc"
	"github.com/scalarorg/protocol-signer/packages/evm"
)

type Handler struct {
	evms   []evm.EvmClient
	signer *btc.PsbtSigner
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

func NewHandler(evms []evm.EvmClient, s *btc.PsbtSigner) (*Handler, error) {
	return &Handler{
		evms:   evms,
		signer: s,
	}, nil
}

func (h *Handler) getEvmClient(chainName string) (*evm.EvmClient, error) {
	for _, evm := range h.evms {
		if evm.ChainName() == chainName {
			return &evm, nil
		}
	}
	return nil, fmt.Errorf("evm client not found for chain name: %s", chainName)
}
