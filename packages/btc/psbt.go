package btc

import (
	"bytes"
	"fmt"

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
func signSegWitV1KeySpend(in *psbt.PInput, tx *wire.MsgTx,
	sigHashes *txscript.TxSigHashes, idx int, privKey *btcec.PrivateKey,
	tapscriptRootHash []byte) error {

	rawSig, err := txscript.RawTxInTaprootSignature(
		tx, sigHashes, idx, in.WitnessUtxo.Value,
		in.WitnessUtxo.PkScript, tapscriptRootHash, in.SighashType,
		privKey,
	)
	if err != nil {
		return fmt.Errorf("error signing taproot input %d: %v", idx,
			err)
	}

	in.TaprootKeySpendSig = rawSig

	return nil
}

// signSegWitV1ScriptSpend attempts to generate a signature for a SegWit version
// 1 (p2tr) input and stores it in the TaprootScriptSpendSig field.
func signSegWitV1ScriptSpend(in *psbt.PInput, tx *wire.MsgTx,
	sigHashes *txscript.TxSigHashes, idx int, privKey *btcec.PrivateKey,
	leaf txscript.TapLeaf) error {

	rawSig, err := txscript.RawTxInTapscriptSignature(
		tx, sigHashes, idx, in.WitnessUtxo.Value,
		in.WitnessUtxo.PkScript, leaf, in.SighashType, privKey,
	)
	if err != nil {
		return fmt.Errorf("error signing taproot script input %d: %v",
			idx, err)
	}

	// Get the 33-byte compressed public key
	compressedPubKey := privKey.PubKey().SerializeCompressed()

	// Extract the last 32 bytes (the x-coordinate) for the x-only public key
	xOnlyPubKey := compressedPubKey[1:] // Skip the first byte (prefix)

	// toXOnlyPubKey from privKey
	in.TaprootBip32Derivation = append(
		in.TaprootBip32Derivation, &psbt.TaprootBip32Derivation{
			XOnlyPubKey: xOnlyPubKey,
		},
	)

	leafHash := leaf.TapHash()

	in.TaprootScriptSpendSig = append(
		in.TaprootScriptSpendSig, &psbt.TaprootScriptSpendSig{
			XOnlyPubKey: in.TaprootBip32Derivation[0].XOnlyPubKey,
			LeafHash:    leafHash[:],
			// We snip off the sighash flag from the end (if it was
			// specified in the first place.)
			Signature: rawSig[:schnorr.SignatureSize],
			SigHash:   in.SighashType,
		},
	)

	return nil
}

// 020000000001015ad096e8fbceb648fd1278877a57298e5441c91bde19b523b87eaf0a94557c090000000000fdffffff01722600000000000016001450dceca158a9c872eb405d52293d351110572c9e0440f3d657e1fb6972b8992df96055d27296fe41a185d34486b77ae6b3fcd87cda31855b7aac430b153dbdb43037911fccf6f010de3b188341dc5aa6aa5e11b2debe40c5bcf4a18cca8651d24ce8ad70c6b4a60dca874a0cae01b4165769556a411dff7dad8bb7e637cadaa9be1477b7b69e799bd1f7b7aab85add3d5b6e69be22ab5444202ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350aad20cf5dff57a173c5ac8323c4baca3fff0728eb716f39f0e5a60312320cd2935b0cac61c150929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0788050e79d530637b2bf963ec79e739ea478978b77b362649439e20045cdcb566ee53347bcebe6c4c52f0b194b8ac3a58febe0d1ac65227c7b4b1420ee4911cc00000000

// 020000000001015ad096e8fbceb648fd1278877a57298e5441c91bde19b523b87eaf0a94557c090000000000fdffffff01722600000000000016001450dceca158a9c872eb405d52293d351110572c9e0440f3d657e1fb6972b8992df96055d27296fe41a185d34486b77ae6b3fcd87cda31855b7aac430b153dbdb43037911fccf6f010de3b188341dc5aa6aa5e11b2debe40c5bcf4a18cca8651d24ce8ad70c6b4a60dca874a0cae01b4165769556a411dff7dad8bb7e637cadaa9be1477b7b69e799bd1f7b7aab85add3d5b6e69be22ab5444202ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350aad20cf5dff57a173c5ac8323c4baca3fff0728eb716f39f0e5a60312320cd2935b0cac61c150929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0788050e79d530637b2bf963ec79e739ea478978b77b362649439e20045cdcb566ee53347bcebe6c4c52f0b194b8ac3a58febe0d1ac65227c7b4b1420ee4911cc00000000
