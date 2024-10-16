package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/scalarorg/protocol-signer/packages/types"
)

// var _ ExternalBtcSigner = (*PsbtSigner)(nil)

type PsbtSigner struct {
	client *BtcClient
}

func NewPsbtSigner(client *BtcClient) *PsbtSigner {
	return &PsbtSigner{
		client: client,
	}
}

func (s *PsbtSigner) SignPsbt(psbtPacket *psbt.Packet) (*types.SigningResult, error) {

	err := s.client.UnlockWallet(60, "passphrase")
	if err != nil {
		return nil, fmt.Errorf("failed to unlock wallet: %w", err)
	}

	signedPacket, err := s.client.SignPsbt(psbtPacket)

	fmt.Printf("signedPacket: %v\n", signedPacket)

	if err != nil {
		return nil, fmt.Errorf("failed to sign PSBT packet: %w", err)
	}

	if len(signedPacket.Inputs[0].TaprootScriptSpendSig) == 0 {
		// this can happen if btcwallet does not maintain the private key for the
		// for the public in signing request
		return nil, fmt.Errorf("no signature found in PSBT packet. Wallet does not maintain covenant public key")
	}

	schnorSignature := signedPacket.Inputs[0].TaprootScriptSpendSig[0].Signature

	parsedSignature, err := schnorr.ParseSignature(schnorSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schnorr signature in psbt packet: %w", err)

	}

	fmt.Println("parsedSignature: ", parsedSignature)

	result := &types.SigningResult{
		Signature: parsedSignature,
	}

	return result, nil
}
