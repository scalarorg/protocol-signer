package handlers_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/scalarorg/protocol-signer/config"
	"github.com/scalarorg/protocol-signer/internals/signer/handlers"
	"github.com/scalarorg/protocol-signer/packages/btc"
	"github.com/scalarorg/protocol-signer/packages/evm"
)

var mockHandler *handlers.Handler

func TestMain(m *testing.M) {

	cfg, err := config.GetConfig("../../../example/config-test.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("cfg: %v\n", cfg)

	parsedConfig, err := cfg.Parse()
	if err != nil {
		panic(err)
	}

	broadcaster, err := btc.NewBtcClient(parsedConfig.BtcNodeConfig)
	if err != nil {
		panic(err)
	}

	signerClient, err := btc.NewBtcClient(parsedConfig.BtcSignerConfig)
	if err != nil {
		panic(err)
	}

	signer := btc.NewPsbtSigner(signerClient, parsedConfig.BtcSignerConfig.Address, parsedConfig.BtcSignerConfig.Passphrase)

	evmClients := make([]evm.EvmClient, len(cfg.EvmConfigs))
	for i, evmConfig := range cfg.EvmConfigs {
		evmClients[i] = *evm.NewEvmClient(evmConfig)
	}
	mockHandler, err = handlers.NewHandler(evmClients, signer, broadcaster, "mock token")
	if err != nil {
		panic(err)
	}
	m.Run()
}

// Note: To run this test, must build bitcoin-vault-ffi first then copy to the lib folder

// CGO_LDFLAGS="-L./lib -lbitcoin_vault_ffi" CGO_CFLAGS="-I./lib" go test -timeout 10m -run ^TestSignUnBondingIngoreCheckOnEvm$ github.com/scalarorg/protocol-signer/internals/signer/handlers -v -count=1
func TestSignUnBondingIngoreCheckOnEvm(t *testing.T) {
	userSignedPsbtBase64 := "cHNidP8BAFICAAAAAd9iI5e5JA2hwzNMPhrUMq9mpr7KmufLlOdDIlIVKGGeAAAAAAD9////ASgjAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKxAnAAAAAAAAIlEg2t54XUPHU7zIxm8h/vBWQ+u02YEqpgeCx0QBiSVbu0sBAwQAAAAAQRQq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1CoshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfQArd2Xo+Hhqn2vZLd+Q1bWKcM5GK/qQMnBAifro0JaJY4ypgbsFBAFqbTbaPXrzcnSxxZdJyJp2BexQzqe140vlCFcBQkpt0waBJVLeLS2A16XpeB4paDyjsltVHv+6azoA6wLA+TxG6WUpeNIqF9MLRbzubGb4+7/SUwSwwUBU1hSVQRSAq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1Cq0gE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dyswCEWE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dwlAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAACEWKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQolAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAAAEXIFCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAARggFq1jps0IRdYmx2Bk96iRBq0h1JkJ5fGwYFdUZNhvlgYAAA=="

	fmt.Printf("userSignedPsbtBase64: %s\n", userSignedPsbtBase64)

	packet, err := psbt.NewFromRawBytes(strings.NewReader(userSignedPsbtBase64), true)
	if err != nil {
		t.Fatalf("Failed to parse PSBT: %v", err)
	}

	fmt.Println("Before signing")

	fmt.Printf("packet: %+v\n", packet.Inputs[0])

	finalTx, err := mockHandler.SignPsbt(packet)
	if err != nil {
		t.Fatalf("Failed to sign PSBT: %v", err)
	}

	for _, txIn := range finalTx.TxIn {
		fmt.Printf("txIn: %+v\n", txIn)
	}

	for _, txOut := range finalTx.TxOut {
		fmt.Printf("txOut: %+v\n", txOut)
	}

	// log witness
	for _, txIn := range finalTx.TxIn {
		fmt.Printf("txIn: %+v\n", txIn.Witness)
	}

	// txid, err := mockHandler.BroadcastTx(finalTx)
	// if err != nil {
	// 	t.Fatalf("Failed to send raw transaction: %v", err)
	// }

	// fmt.Printf("txid: %s\n", txid)
}
