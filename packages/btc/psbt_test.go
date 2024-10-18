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
		Host:    "192.168.1.34:18332",
		User:    "user",
		Pass:    "password",
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
	btcAddressString := "bcrt1qp0pyg95qusxrvqn2vggjz926u0gvxk8psuwpw6"

	btcAddress, err := btcutil.DecodeAddress(btcAddressString, &chaincfg.RegressionNetParams)
	if err != nil {
		t.Fatalf("Failed to decode address: %v", err)
	}

	privKey, err := signerClient.DumpPrivateKey(btcAddress)
	if err != nil {
		t.Fatalf("Failed to dump private key: %v", err)
	}

	fmt.Println("privKey", privKey)

	stakerPsbt := "cHNidP8BAFICAAAAAfplYlPYiWq48xv17k+OjtsWPWOJ2t1Iwv9RMlpBOSSYAAAAAAD9////Adi9AAAAAAAAFgAUmZG33OeZPn7n2tGSGWLWM2jk2ywAAAAAAAEBK4A4AQAAAAAAIlEgZoZJC6/8BXk60MZskWty6TNnjVmiQmIf4Higz2oTwLxBFE2u4mlM/5UIwIvQi7E9BvJ82TRF2eUzBgjd0HLFOZdpyZwNkOxbWoGdiVt2TnXUBrJA/oLeJ9lEgzqUyEQQ/9ZANKTY4jop4juxWM0Hm4vt2yMzgha8UnjIWsp+txlEbL+QC64Cegzedo92/YLf2xtnFQJZrpz8g/1oOxwY+bKzuWIVwVCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrA8bHBS6OxMxFOaySwPemFsAipfLUVGN+1L5k4xErTLtVlsDbZjCe3tS9xwW+YzniA+XD7CDZJaOPE5ZxucesTIEUgTa7iaUz/lQjAi9CLsT0G8nzZNEXZ5TMGCN3QcsU5l2mtIBC+QfJMDdXQjFhk6kCYegEgCgTq56IKOxkia2W6PMDerMAAAA=="

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
