package handlers_test

import (
	"bytes"
	"encoding/hex"
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
// cp -p ../../bitcoin-vault/target/release/libbitcoin_vault_ffi.* ./lib

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

	buf := new(bytes.Buffer)
	finalTx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())

	fmt.Printf("finalTx: %s\n", txHex)

	if txHex != "02000000000101df622397b9240da1c3334c3e1ad432af66a6beca9ae7cb94e74322521528619e0000000000fdffffff01282300000000000016001450dceca158a9c872eb405d52293d351110572c9e0440e7757536ce5cf4246485d74cfc14d74821acefd8ccba32c92a1c6852b1219d82320ea9868bb57e552be3bc8fea5f9a214006d4d2e87476df815bce6a18f68d4e400addd97a3e1e1aa7daf64b77e4356d629c33918afea40c9c10227eba3425a258e32a606ec141005a9b4db68f5ebcdc9d2c7165d272269d817b1433a9ed78d2f944202ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350aad201387aab21303782b17e760c670432559df3968e52cb82cc2d8f9be43a227d5dcac41c050929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0b03e4f11ba594a5e348a85f4c2d16f3b9b19be3eeff494c12c3050153585255000000000" {
		t.Fatalf("txHex is empty")
	}



	// txid, err := mockHandler.BroadcastTx(finalTx)
	// if err != nil {
	// 	t.Fatalf("Failed to send raw transaction: %v", err)
	// }

	// fmt.Printf("txid: %s\n", txid)
}
