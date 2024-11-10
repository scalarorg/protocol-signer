package types

import (
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// SignUnbondingTxPayload carries all data necessary to sign unbonding transaction
// TODO: add validations
type SignUnbondingTxRequest struct {
	EvmChainName        string `json:"evm_chain_name"`
	EvmTxID             string `json:"evm_tx_id"`
	UnbondingPsbtBase64 string `json:"unbonding_psbt_base64"` // base64 encoded psbt, if you want to change raw hex, please think about your life.
}

// SignUnbondingTxResponse covenant member schnorr signature
type SignAndBroadcastPsbtReponse struct {
	TxID  *chainhash.Hash `json:"tx_id"`
	TxHex string          `json:"tx_hex"`
}

type SigningResult struct {
	Signature *schnorr.Signature
}
