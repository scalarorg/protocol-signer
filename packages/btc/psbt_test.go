package btc_test

import (
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

	signerClient.UnlockWallet(60, "passphrase")

	m.Run()
}

func TestSignPsbt(t *testing.T) {
	wifKey := "cURd2WaKTJsyHgRM5kViqXvjo4R3oyT1gfV5CJhQxD1EuprcnSSP"

	wif, err := btcutil.DecodeWIF(wifKey)
	if err != nil {
		t.Fatalf("Failed to decode WIF: %v", err)
	}

	privKey := wif.PrivKey

	stakerPsbt := "cHNidP8BAFICAAAAAVrQluj7zrZI/RJ4h3pXKY5UQckb3hm1I7h+rwqUVXwJAAAAAAD9////AXImAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKw8nAAAAAAAAIlEgkFj0s5tPU4QO+QWJFepTzvTiZxWVJMmxkN+zd/AYnZtBFCrjHqhwmu2oGUuj4vfn6V5oDotlE1yJg8CimNF7xTUKa/PeGFPFI+tlh34WOxsYSsPX8wg4uJXl35HhYElEst9A89ZX4ftpcriZLflgVdJylv5BoYXTRIa3euaz/Nh82jGFW3qsQwsVPb20MDeRH8z28BDeOxiDQdxapqpeEbLevmIVwVCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAeIBQ551TBjeyv5Y+x55znqR4l4t3s2JklDniAEXNy1Zu5TNHvOvmxMUvCxlLisOlj+vg0axlInx7SxQg7kkRzEUgKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQqtIM9d/1ehc8WsgyPEuso//wco63FvOfDlpgMSMgzSk1sMrMAAAA=="

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

	fmt.Println("result: ", result)
}
