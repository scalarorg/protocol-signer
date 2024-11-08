package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	psbtFfi "github.com/scalarorg/bitcoin-vault/ffi/go-psbt"
	"github.com/scalarorg/protocol-signer/packages/types"
)

func (s *PsbtSigner) SignPsbt(psbtPacket *psbt.Packet) (*wire.MsgTx, error) {
	//TODO: fix hardcode
	err := s.client.UnlockWallet(60, s.passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock wallet: %w", err)
	}

	fmt.Printf("s: %+v\n", s)

	privKey, err := s.client.DumpPrivateKey(s.address)
	if err != nil {
		return nil, fmt.Errorf("failed to dump private key: %w", err)
	}

	fmt.Printf("privKey: %+v\n", privKey)

	fmt.Printf("privBytes: %x\n", privKey.Serialize())

	// 2a8721658a12c63f4aeb5548f4988c25842602c6564303cede678d9a92178ef0

	privKeyBytes, _ := hex.DecodeString("f5b5ce21907a33c4b39d50649bcbc7ee029a3905c6ee470e7b434fbc960c794a")

	var buf bytes.Buffer
	err = psbtPacket.Serialize(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize psbt: %w", err)
	}
	psbtBytes := buf.Bytes()

	isTestnet := true // TODO: FIX THIS HARDCODE by checking the network

	networkKind := psbtFfi.NetworkKindTestnet
	if !isTestnet {
		networkKind = psbtFfi.NetworkKindMainnet
	}

	isFinalized := true

	tx, err := psbtFfi.SignPsbtBySingleKey(
		psbtBytes,       // []byte containing PSBT
		privKeyBytes[:], // []byte containing private key
		networkKind,     // bool indicating if testnet
		isFinalized,     // finalize
	)
	if err != nil {
		log.Fatal(err)
	}

	finalTx := &wire.MsgTx{}
	err = finalTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		return nil, types.NewErrorWithMsg(http.StatusInternalServerError, types.InternalServiceError, err.Error())
	}

	return finalTx, nil
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

func FinalizePsbt(packet *psbt.Packet, signedInputs []uint32, sortSigs bool) error {
	for idx := range signedInputs {
		input := &packet.Inputs[idx]
		if sortSigs {
			leaf := input.TaprootLeafScript[0]
			convertedLeaf := txscript.TapLeaf{
				LeafVersion: leaf.LeafVersion,
				Script:      leaf.Script,
			}
			sortTaprootSigsByPubKey(input, convertedLeaf)
		}

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

// func sortTaprootSigsByPubKey(input *psbt.PInput) error {
// 	// revert the order of the input.TaprootScriptSpendSig
// 	sort.Slice(input.TaprootScriptSpendSig, func(i, j int) bool {
// 		return i > j
// 	})
// 	return nil
// }

type taprootScriptSpendSigWithPosition struct {
	*psbt.TaprootScriptSpendSig
	Position int
}

func sortTaprootSigsByPubKey(input *psbt.PInput, leaf txscript.TapLeaf) error {
	leafHash := leaf.TapHash()

	// Filter signatures by leaf hash and add position information
	var sigs []*taprootScriptSpendSigWithPosition
	for _, sig := range input.TaprootScriptSpendSig {
		if bytes.Equal(sig.LeafHash, leafHash[:]) {
			pos, err := pubkeyPositionInScript(sig.XOnlyPubKey, leaf.Script)
			if err != nil {
				return fmt.Errorf("error finding pubkey position: %w", err)
			}
			sigs = append(sigs, &taprootScriptSpendSigWithPosition{
				TaprootScriptSpendSig: sig,
				Position:              pos,
			})
		}
	}

	// Sort signatures by position in descending order
	sort.Slice(sigs, func(i, j int) bool {
		return sigs[i].Position > sigs[j].Position
	})

	// Update the input with sorted signatures
	input.TaprootScriptSpendSig = make([]*psbt.TaprootScriptSpendSig, len(sigs))
	for i, sig := range sigs {
		input.TaprootScriptSpendSig[i] = sig.TaprootScriptSpendSig
	}

	return nil
}

func pubkeyPositionInScript(pubkey []byte, script []byte) (int, error) {
	pubkeyHash := btcutil.Hash160(pubkey)
	pubkeyXOnly := pubkey // Note: pubkey is already x-only in Go implementation

	decompiled, err := txscript.DisasmString(script)
	if err != nil {
		return -1, fmt.Errorf("error decompiling script: %w", err)
	}

	elements := strings.Split(decompiled, " ")
	for i, element := range elements {
		data, err := hex.DecodeString(element)
		if err != nil {
			continue // Skip non-hex elements (opcodes)
		}
		if bytes.Equal(data, pubkey) || bytes.Equal(data, pubkeyHash) || bytes.Equal(data, pubkeyXOnly) {
			return i, nil
		}
	}

	return -1, fmt.Errorf("pubkey not found in script")
}
