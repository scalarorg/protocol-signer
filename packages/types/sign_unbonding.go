package types

import (
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// SignUnbondingTxPayload carries all data necessary to sign unbonding transaction
// TODO: add validations
type SignUnbondingTxRequest struct {
	EvmChainName  string `json:"evm_chain_name"`
	EvmTxID       string `json:"evm_tx_id"`
	UnbondingPsbt string `json:"unbonding_psbt"` // base64 encoded psbt
}

// SignUnbondingTxResponse covenant member schnorr signature
type SignUnbondingTxResponse struct {
	SignatureHex string `json:"signature_hex"`
}

type SigningResult struct {
	Signature *schnorr.Signature
}
