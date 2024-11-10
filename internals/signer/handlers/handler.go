package handlers

import (
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/scalarorg/protocol-signer/packages/btc"
	"github.com/scalarorg/protocol-signer/packages/evm"
)

type Handler struct {
	evms        []evm.EvmClient
	signer      *btc.PsbtSigner
	broadcaster btc.BtcClientInterface
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

func NewHandler(evms []evm.EvmClient, s *btc.PsbtSigner, b btc.BtcClientInterface, t string) (*Handler, error) {
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

func (h *Handler) SignPsbt(packet *psbt.Packet) (*wire.MsgTx, error) {
	return h.signer.SignPsbt(packet)
}

func (h *Handler) BroadcastTx(tx *wire.MsgTx) (*chainhash.Hash, error) {
	return h.broadcaster.SendTx(tx)
}

func (h *Handler) TestMempoolAccept(txs []*wire.MsgTx, maxFeeRatePerKb float64) ([]*btcjson.TestMempoolAcceptResult, error) {
	return h.broadcaster.TestMempoolAccept(txs, maxFeeRatePerKb)
}
