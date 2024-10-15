package types

// SignUnbondingTxPayload carries all data necessary to sign unbonding transaction
type SignUnbondingTxRequest struct {
	//Souce evm chain name for evm call to check if unbonding is valid on the source chain
	ChainName string `json:"chain_name"`
	// TxId is the transaction id of the unbonding transaction on the source chain
	TxId                     string `json:"tx_id"`
	StakingOutputPkScriptHex string `json:"staking_output_pk_script_hex"`
	UnbondingTxHex           string `json:"unbonding_tx_hex"`
	StakerUnbondingSigHex    string `json:"staker_unbonding_sig_hex"`
	// 33 bytes compressed public key
	CovenantPublicKey string `json:"covenant_public_key"`
}

// SignUnbondingTxResponse covenant member schnorr signature
type SignUnbondingTxResponse struct {
	SignatureHex string `json:"signature_hex"`
}
