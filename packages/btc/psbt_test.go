package btc_test

// import (
// 	"bytes"
// 	"fmt"
// 	"strings"
// 	"testing"

// 	"github.com/btcsuite/btcd/btcec/v2"
// 	"github.com/btcsuite/btcd/btcutil"
// 	"github.com/btcsuite/btcd/btcutil/psbt"
// 	"github.com/btcsuite/btcd/chaincfg"
// 	"github.com/scalarorg/protocol-signer/config"
// 	"github.com/scalarorg/protocol-signer/packages/btc"
// )

// var broadcasterClient *btc.BtcClient
// var psbtSigner *btc.PsbtSigner

// func TestMain(m *testing.M) {
// 	c, err := btc.NewBtcClient(&config.ParsedBtcConfig{
// 		Host:    "testnet3.btc.scalar.org:80",
// 		User:    "mike",
// 		Pass:    "apd3g41pkl",
// 		Network: &chaincfg.TestNet3Params,
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	broadcasterClient = c

// 	broadcasterClient.UnlockWallet(60, "passphrase")

// 	m.Run()
// }

// // go test -run ^TestSignPsbt$ github.com/scalarorg/protocol-signer/packages/btc -v -count=1

// // Deprecated: Currently, use the ffi/go-psbt to sign the PSBT
// func TestSignPsbt(t *testing.T) {
// 	// btcAddressString := "tb1qwdw0aymdetu7gnfajk39gfx6w5wtavggzrv6nd"

// 	// btcAddress, err := btcutil.DecodeAddress(btcAddressString, &chaincfg.RegressionNetParams)
// 	// if err != nil {
// 	// 	t.Fatalf("Failed to decode address: %v", err)
// 	// }

// 	wif, err := btcutil.DecodeWIF("cNGbmJbymnzaFUPZ8XSLvQQxHEEcTkh1ojBMMpvg5vFX5V1afcmR")
// 	if err != nil {
// 		t.Fatalf("Failed to decode WIF: %v", err)
// 	}

// 	privKey := wif.PrivKey
// 	pubKey := privKey.PubKey()
// 	fmt.Printf("pubKey: %x\n", pubKey.SerializeCompressed())
// 	btcPubKey := btcec.PublicKey(*pubKey)
// 	btcAdress, _ := btcutil.NewAddressPubKey(btcPubKey.SerializeCompressed(), &chaincfg.TestNet3Params)

// 	psbtSigner = btc.NewPsbtSigner(broadcasterClient, btcAdress.String(), "passphrase")

// 	fmt.Println("privKey", privKey)

// 	stakerPsbt := "cHNidP8BAFICAAAAATNTX+s8B0R4iNMYYnfz7y4pgry7ZymAMFwCcjM7Tf8EAAAAAAD9////AQc5AQAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBK/BJAgAAAAAAIlEgUiPouA919Arm/mld5UVulY4HFUe3CUvidmtDCJyUOHRBFCrjHqhwmu2oGUuj4vfn6V5oDotlE1yJg8CimNF7xTUKGqpwf1oaV/8WEF+uunL1V6iWOpitlLONjh+Oud3IvilAVhiwkYUAI1N4oBZBVgPv08pC3oYssJivGLqRgu/anRBW0jxBMZmMykQpdhzAPd7dZVcxuIegs9EiXPmqY/gWc2IVwVCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAFuIuioNMBacc5AXCkPpmmlMAQJR7esmfH/PGLWmWX25mR0dHHl1WKnImL8iqD3X/VURzLq83wrCHDtY3AZcpEkUgKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQqtIPNYZ068zmiDKsMCuC7Hzk6c3Y+FZoTKPVRWhqJAvhXYrMAAAA=="

// 	packet, err := psbt.NewFromRawBytes(strings.NewReader(stakerPsbt), true)
// 	if err != nil {
// 		t.Fatalf("Failed to parse PSBT: %v", err)
// 	}

// 	finalTx, err := psbtSigner.SignPsbt(packet, false)
// 	if err != nil {
// 		t.Fatalf("Failed to sign PSBT: %v", err)
// 	}

// 	var txBytes bytes.Buffer
// 	finalTx.Serialize(&txBytes)

// 	fmt.Printf("TxHex1: %x\n", txBytes.Bytes())

// 	txid, err := broadcasterClient.RpcClient.SendRawTransaction(finalTx, false)
// 	if err != nil {
// 		t.Logf("Failed to send raw transaction: %v", err)
// 		newPacket, err := psbt.NewFromRawBytes(strings.NewReader(stakerPsbt), true)
// 		if err != nil {
// 			t.Fatalf("Failed to parse PSBT: %v", err)
// 		}

// 		finalTx, err = psbtSigner.SignPsbt(newPacket, true)
// 		if err != nil {
// 			t.Fatalf("Failed to sign PSBT: %v", err)
// 		}

// 		var txBytes bytes.Buffer
// 		finalTx.Serialize(&txBytes)

// 		fmt.Printf("TxHex2: %x\n", txBytes.Bytes())

// 		txid, err = broadcasterClient.RpcClient.SendRawTransaction(finalTx, false)
// 		if err != nil {
// 			t.Fatalf("Failed to send raw transaction: %v", err)
// 		}
// 	}

// 	fmt.Println("txid", txid)
// }
