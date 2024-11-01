package handlers

import (
	"fmt"
	"net/http"

	"github.com/scalarorg/protocol-signer/packages/btc"
	"github.com/scalarorg/protocol-signer/packages/evm"
)

type Handler struct {
	evms        []evm.EvmClient
	signer      *btc.PsbtSigner
	broadcaster *btc.BtcClient
	token       string
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

func NewHandler(evms []evm.EvmClient, s *btc.PsbtSigner, b *btc.BtcClient, t string) (*Handler, error) {
	if len(evms) == 0 {
		return nil, fmt.Errorf("no evm clients provided")
	}

	if s == nil {
		return nil, fmt.Errorf("no btc signer provided")
	}

	if b == nil {
		return nil, fmt.Errorf("no btc broadcaster provided")
	}

	if t == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	return &Handler{
		evms:        evms,
		signer:      s,
		broadcaster: b,
		token:       t,
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
