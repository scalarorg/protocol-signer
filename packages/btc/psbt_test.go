package btc_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

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

	// signerClient.UnlockWallet(60, "passphrase")

	m.Run()
}

// go test -run ^TestSignPsbt$ github.com/scalarorg/protocol-signer/packages/btc -v -count=1

func TestSignPsbt(t *testing.T) {
	wifKey := "cU5n3yJpmFjEc7oDAE8ceF3Xr6Gknhg7JnJ8XABcoSYFttfwwoob"

	wif, err := btcutil.DecodeWIF(wifKey)
	if err != nil {
		t.Fatalf("Failed to decode WIF: %v", err)
	}

	privKey := wif.PrivKey

	stakerPsbt := "cHNidP8BAFICAAAAAX6WwhSPsG0P3RAx5xkdyDWJlkYnN8Bfj0I1+5RamcYDAAAAAAD9////Adi9AAAAAAAAFgAUN5/Vi2cSjSWviTHiiHPNBhTOxDUAAAAAAAEBK4A4AQAAAAAAIlEgg8ldB/Oq2r58Txahcvukuw/cv/fFmrQQkgSY4uUBdr1BFP3Uoq41F+2m2tGmjm9ZyeelRV9bSx80dYaA0FKPNhPcF5d5DZaJDtibdN+87atu1GbbBvxj5kgBue9X4xx1rS5AwjqZRntI2P0WpqT5J+HVV+Fo6hBZ6Gv8MMHQyqfNSiQB+Zc8CLY62HW96ex3jtZX9Zicb9Kytd3QEloUVjgATGIVwVCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAPu8mVK0CpF3NA1w0ywWJmEbtIwnBNVFh4KTgYZ9lA8z1x6NlDe6V7Dej0LxpwcCSXOrYRIpOnuhwZQKLoZQkPUUg/dSirjUX7aba0aaOb1nJ56VFX1tLHzR1hoDQUo82E9ytILmYDQKMuOD7t4mD+hfxHzMVJrxwgNWAa3am3/F2mshmrMAAAA=="

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

	fmt.Println("finalTx", hex.EncodeToString(buf.Bytes()))
}
