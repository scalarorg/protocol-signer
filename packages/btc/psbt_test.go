package btc_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/scalarorg/protocol-signer/config"
	"github.com/scalarorg/protocol-signer/packages/btc"
)

var signerClient *btc.BtcClient

func TestMain(m *testing.M) {
	c, err := btc.NewBtcClient(&config.ParsedBtcConfig{
		Host:    "testnet3.btc.scalar.org:80",
		User:    "mike",
		Pass:    "apd3g41pkl",
		Network: &chaincfg.TestNet3Params,
	})

	if err != nil {
		panic(err)
	}

	signerClient = c

	signerClient.UnlockWallet(60, "passphrase")

	m.Run()
}

// go test -run ^TestSignPsbt$ github.com/scalarorg/protocol-signer/packages/btc -v -count=1

func TestSignPsbt(t *testing.T) {
	btcAddressString := "tb1qwdw0aymdetu7gnfajk39gfx6w5wtavggzrv6nd"

	btcAddress, err := btcutil.DecodeAddress(btcAddressString, &chaincfg.RegressionNetParams)
	if err != nil {
		t.Fatalf("Failed to decode address: %v", err)
	}

	privKey, err := signerClient.DumpPrivateKey(btcAddress)
	if err != nil {
		t.Fatalf("Failed to dump private key: %v", err)
	}

	fmt.Println("privKey", privKey)

	stakerPsbt := "cHNidP8BAFICAAAAAflAdaiII1fkBTGGT8vk5a3YUWMTPqgFB/KD0IpN9O/aAAAAAAD9////AXjiCgAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKwA1DAAAAAAAIlEgUiPouA919Arm/mld5UVulY4HFUe3CUvidmtDCJyUOHRBFCrjHqhwmu2oGUuj4vfn6V5oDotlE1yJg8CimNF7xTUKGqpwf1oaV/8WEF+uunL1V6iWOpitlLONjh+Oud3IvilAW9mng7KEV1Vas9AVDCPPKGEDxkZD7gml1V5/7wS3xuOsiDffs8Gi4zChnzTWEx6q5/M/Yl/lbvunNGt2AdpVu2IVwVCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAFuIuioNMBacc5AXCkPpmmlMAQJR7esmfH/PGLWmWX25mR0dHHl1WKnImL8iqD3X/VURzLq83wrCHDtY3AZcpEkUgKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQqtIPNYZ068zmiDKsMCuC7Hzk6c3Y+FZoTKPVRWhqJAvhXYrMAAAA=="

	packet, err := psbt.NewFromRawBytes(strings.NewReader(stakerPsbt), true)
	if err != nil {
		t.Fatalf("Unable to parse Psbt: %v", err)
	}

	result, err := btc.SignPsbtAll(packet, privKey)
	if err != nil {
		t.Fatalf("Failed to sign PSBT: %v", err)
	}

	if result == nil {
		t.Fatalf("No result returned from signing")
	}

	// Proceed with finalization
	err = btc.FinalizePsbt(packet, result)
	if err != nil {
		t.Fatalf("Failed to finalize PSBT: %v", err)
	}

	finalTx, err := psbt.Extract(packet)
	if err != nil {
		t.Fatalf("Failed to extract transaction: %v", err)
	}

	var buf bytes.Buffer
	finalTx.Serialize(&buf)

	txHex := hex.EncodeToString(buf.Bytes())
	fmt.Println("txHex", txHex)

	// Broadcast the transaction

	txid, err := signerClient.RpcClient.SendRawTransaction(finalTx, false)
	if err != nil {
		t.Fatalf("Failed to broadcast transaction: %v", err)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("txid", txid)
}
