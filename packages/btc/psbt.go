package btc

import (
	"bytes"
	"fmt"
	"math"
	"sort"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/scalarorg/protocol-signer/packages/types"
)

func (s *PsbtSigner) SignPsbt(psbtPacket *psbt.Packet) (*types.SigningResult, error) {
	//TODO: fix hardcode
	err := s.client.UnlockWallet(60, "passphrase")
	if err != nil {
		return nil, fmt.Errorf("failed to unlock wallet: %w", err)
	}

	// psbtEncoded, err := psbtPacket.B64Encode()
	// if err != nil {
	// 	return nil, err
	// }

	// {
	// 	result = ""

	// 	decodedBytes, err := base64.StdEncoding.DecodeString(result.Psbt)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	decodedBytes := []byte{}

	signedPacket, err := psbt.NewFromRawBytes(bytes.NewReader(decodedBytes), false)
	if err != nil {
		return nil, err
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

	result := &types.SigningResult{
		Signature: parsedSignature,
	}

	return result, nil
}

// Ref: https://github.com/lightningnetwork/lnd/blob/4da26fb65a669fbee68fa36e60259a8da8ef6d3b/lnwallet/btcwallet/psbt.go#L132

func SignPsbtAll(packet *psbt.Packet, privKey *secp256k1.PrivateKey) ([]uint32, error) {
	var signedInputs []uint32

	tx := packet.UnsignedTx
	prevOutputFetcher := wallet.PsbtPrevOutputFetcher(packet)
	sigHashes := txscript.NewTxSigHashes(tx, prevOutputFetcher)

	for idx := range packet.Inputs {

		input := &packet.Inputs[idx]

		if input.WitnessUtxo == nil {
			continue
		}

		if len(input.FinalScriptWitness) > 0 {
			continue
		}

		// Schnorr key path signature (Taproot key spend)
		// if input.TaprootMerkleRoot == nil {
		// 	rootHash := make([]byte, 0)
		// 	err := signSegWitV1KeySpend(&input, tx, sigHashes, idx, privKey, rootHash)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// } else {
		// Schnorr + script path (Taproot script spend)
		leafScript := input.TaprootLeafScript[0]
		leaf := txscript.TapLeaf{
			LeafVersion: leafScript.LeafVersion,
			Script:      leafScript.Script,
		}
		err := signSegWitV1ScriptSpend(input, tx, sigHashes, idx, privKey, leaf)
		if err != nil {
			return nil, err
		}

		signedInputs = append(signedInputs, uint32(idx))
	}

	return signedInputs, nil
}

// signSegWitV1KeySpend attempts to generate a signature for a SegWit version 1
// (p2tr) input and stores it in the TaprootKeySpendSig field.
// func signSegWitV1KeySpend(in *psbt.PInput, tx *wire.MsgTx,
// 	sigHashes *txscript.TxSigHashes, idx int, privKey *btcec.PrivateKey,
// 	tapscriptRootHash []byte) error {

// 	rawSig, err := txscript.RawTxInTaprootSignature(
// 		tx, sigHashes, idx, in.WitnessUtxo.Value,
// 		in.WitnessUtxo.PkScript, tapscriptRootHash, in.SighashType,
// 		privKey,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error signing taproot input %d: %v", idx,
// 			err)
// 	}

// 	in.TaprootKeySpendSig = rawSig

// 	return nil
// }

// signSegWitV1ScriptSpend attempts to generate a signature for a SegWit version
// 1 (p2tr) input and stores it in the TaprootScriptSpendSig field.
func signSegWitV1ScriptSpend(in *psbt.PInput, tx *wire.MsgTx,
	sigHashes *txscript.TxSigHashes, idx int, privKey *btcec.PrivateKey,
	leaf txscript.TapLeaf) error {

	inputFetcher := txscript.NewCannedPrevOutputFetcher(
		in.WitnessUtxo.PkScript, in.WitnessUtxo.Value,
	)
	hType := txscript.SigHashDefault

	tapLeafHash := leaf.TapHash()

	sigHash, err := txscript.CalcTapscriptSignaturehash(sigHashes, hType, tx, idx, inputFetcher, leaf, txscript.WithBaseTapscriptVersion(math.MaxUint32, tapLeafHash[:]))
	if err != nil {
		return err
	}

	zeroBytes := [32]byte{}
	signature, err := schnorr.Sign(privKey, sigHash, schnorr.CustomNonce(zeroBytes))
	if err != nil {
		return err
	}

	rawSig := signature.Serialize()

	// Finally, append the sighash type to the final sig if it's not the
	// default sighash value (in which case appending it is disallowed).
	if hType != txscript.SigHashDefault {
		rawSig = append(rawSig, byte(hType))
	}

	XOnlyPubkey := privKey.PubKey().SerializeCompressed()[1:]

	leafHash := leaf.TapHash()

	in.TaprootScriptSpendSig = append(
		in.TaprootScriptSpendSig, &psbt.TaprootScriptSpendSig{
			XOnlyPubKey: XOnlyPubkey,
			LeafHash:    leafHash[:],
			// We snip off the sighash flag from the end (if it was
			// specified in the first place.)
			Signature: rawSig[:schnorr.SignatureSize],
			SigHash:   in.SighashType,
		},
	)

	return nil
}

func FinalizePsbt(packet *psbt.Packet, signedInputs []uint32) error {
	for idx := range signedInputs {
		input := &packet.Inputs[idx]
		sortTaprootSigsByPubKey(input)
		success, err := psbt.MaybeFinalize(packet, idx)
		if err != nil {
			return err
		}
		if !success {
			return fmt.Errorf("failed to finalize PSBT")
		}
	}
	return nil
}

func sortTaprootSigsByPubKey(input *psbt.PInput) {
	sort.Slice(input.TaprootScriptSpendSig, func(i, j int) bool {
		return bytes.Compare(input.TaprootScriptSpendSig[i].XOnlyPubKey[:],
			input.TaprootScriptSpendSig[j].XOnlyPubKey[:]) < 0
	})
}
