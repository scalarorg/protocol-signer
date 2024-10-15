package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/protocol-signer/signerapp"
	"github.com/scalarorg/protocol-signer/signerservice/types"
	"github.com/scalarorg/protocol-signer/utils"
)

func parseSchnorrSigFromHex(hexStr string) (*schnorr.Signature, error) {
	sigBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	return schnorr.ParseSignature(sigBytes)
}

func (h *Handler) SignUnbonding(request *http.Request) (*Result, *types.Error) {
	// Check Authorization header
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, types.NewErrorWithMsg(http.StatusUnauthorized, types.Forbidden, "missing Authorization header")
	}

	// Extract the access token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, types.NewErrorWithMsg(http.StatusUnauthorized, types.Forbidden, "invalid Authorization header format")
	}
	accessToken := parts[1]

	// Verify the access token
	if !h.verifyAccessToken(accessToken) {
		return nil, types.NewErrorWithMsg(http.StatusUnauthorized, types.Forbidden, "invalid access token")
	}

	payload := &types.SignUnbondingTxRequest{}
	err := json.NewDecoder(request.Body).Decode(payload)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid request payload")
	}

	pkScript, err := hex.DecodeString(payload.StakingOutputPkScriptHex)

	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid staking output pk script")
	}

	covenantPublicKeyBytes, err := hex.DecodeString(payload.CovenantPublicKey)

	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid covenant public key")
	}

	covenantPublicKey, err := btcec.ParsePubKey(covenantPublicKeyBytes)

	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid covenant public key")
	}

	unbondingTx, _, err := utils.NewBTCTxFromHex(payload.UnbondingTxHex)

	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid unbonding transaction")
	}

	stakerUnbondingSig, err := parseSchnorrSigFromHex(payload.StakerUnbondingSigHex)

	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid staker unbonding signature")
	}

	// Check if unbonding is valid on the source chain
	evmClient, err := h.getEvmClient(payload.ChainName)

	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "Chain not found")
	}
	txHash := common.HexToHash(payload.TxId)
	err = evmClient.CheckUnbondingTx(request.Context(), txHash, &payload.UnbondingTxHex)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "Unbonding transaction is not valid on the source chain")
	}
	// do not count the requests with invalid arguments
	h.m.IncReceivedSigningRequests()

	sig, err := h.s.SignUnbondingTransaction(
		request.Context(),
		pkScript,
		unbondingTx,
		stakerUnbondingSig,
		covenantPublicKey,
	)

	if err != nil {
		h.m.IncFailedSigningRequests()

		if errors.Is(err, signerapp.ErrInvalidSigningRequest) {
			return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, err.Error())
		}

		// if this is unknown error, return internal server error
		return nil, types.NewErrorWithMsg(http.StatusInternalServerError, types.InternalServiceError, err.Error())
	}

	resp := types.SignUnbondingTxResponse{
		SignatureHex: hex.EncodeToString(sig.Serialize()),
	}

	h.m.IncSuccessfulSigningRequests()

	return NewResult(resp), nil
}

// verifyAccessToken checks if the provided access token is valid
func (h *Handler) verifyAccessToken(token string) bool {
	// Implement your token verification logic here
	// This could involve checking against a database, calling an authentication service, etc.
	// For this example, we'll just check if the token is not empty
	return token != h.t
}
