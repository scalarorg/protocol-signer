package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/protocol-signer/packages/types"
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

	txHash := common.HexToHash(payload.EvmTxID)

	// Check if unbonding is valid on the source chain
	evmClient, err := h.getEvmClient(payload.EvmChainName)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "Chain not found")
	}

	err = evmClient.CheckUnbondingTx(request.Context(), txHash, payload.UnbondingPsbt)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest,
			fmt.Sprintf("Error checking unbonding tx: %s", err.Error()))
	}

	psbtPacket, err := psbt.NewFromRawBytes(strings.NewReader(payload.UnbondingPsbt), true)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusBadRequest, types.BadRequest, "invalid psbt packet")
	}

	encoded, _ := psbtPacket.B64Encode()

	fmt.Println("Encoded: ", encoded)

	result, err := h.signer.SignPsbt(psbtPacket)
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusInternalServerError, types.InternalServiceError, err.Error())
	}

	return NewResult(result), nil
}

// verifyAccessToken checks if the provided access token is valid
func (h *Handler) verifyAccessToken(token string) bool {
	// Implement your token verification logic here
	// This could involve checking against a database, calling an authentication service, etc.
	// For this example, we'll just check if the token is not empty
	// return token != h.t
	return true
}