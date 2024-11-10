package handlers_test

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"github.com/scalarorg/protocol-signer/config"
	"github.com/scalarorg/protocol-signer/internals/signer/handlers"
	"github.com/scalarorg/protocol-signer/packages/btc"
	"github.com/scalarorg/protocol-signer/packages/evm"
	"github.com/stretchr/testify/assert"
)

var mockHandler *handlers.Handler

const (
	UNSTAKE_54167_PSBT_B64         = "cHNidP8BAFICAAAAASm9dsIF2Yr1NotB752ZgynoNVp0vqnnPtyODFPKKsQkAAAAAAD9////ASgjAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKxAnAAAAAAAAIlEgZ7/zV3gKk4JqREZGrsaBxP8fQxYkRHjA1hH5GnXJO4oBAwQAAAAAQRQq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1CoshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfQIjBR8YVwGp1QTMOJcVadbaJIYPd/gzsGaZ63jci5GVgdn3YMCvZW7K654Xwc0IOTbWaRxj6zvCrzHLehcy5nT5CFcFQkpt0waBJVLeLS2A16XpeB4paDyjsltVHv+6azoA6wAY8WLORYd6jGMAq4zgcTd/6BArojm/prhViwo8tsQKFRSAq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1Cq0gE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dyswCEWE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dwlAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAACEWKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQolAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAAAEXIFCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAARggR4Li5f/hJviWsPse5R7SzU/wp7r8u4szV3KnW5FahpAAAA=="
	UNSTAKE_54167_PSBT_FINAL_TXHEX = "0200000000010129bd76c205d98af5368b41ef9d998329e8355a74bea9e73edc8e0c53ca2ac4240000000000fdffffff01282300000000000016001450dceca158a9c872eb405d52293d351110572c9e0440ec06f885e6f3b24922500c2e70011297900020ed3e5dbfc633fc544353c5b915b7e28acee67f1204f9866706f20275460d963c005557b060b66a8d2961447c334088c147c615c06a7541330e25c55a75b6892183ddfe0cec19a67ade3722e46560767dd8302bd95bb2bae785f073420e4db59a4718facef0abcc72de85ccb99d3e44202ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350aad201387aab21303782b17e760c670432559df3968e52cb82cc2d8f9be43a227d5dcac41c150929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0063c58b39161dea318c02ae3381c4ddffa040ae88e6fe9ae1562c28f2db1028500000000"
	TEST_54155_PSBT_B64            = "cHNidP8BAFICAAAAAYP9BIKyx2PMWhFuMAT3wcODSvWRakPI1WvHVQEmwAdlAAAAAAD9////ASgjAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKxAnAAAAAAAAIlEgZ7/zV3gKk4JqREZGrsaBxP8fQxYkRHjA1hH5GnXJO4oBAwQAAAAAQRQq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1CoshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfQPSfuF9Tags3JLQ1vK1smKCghwYGcWR4AAALh/mb6yNFJ3EdqaUmI3pWW93xLaX+HqUyUEKNKG0DM58FY6aAyxJCFcFQkpt0waBJVLeLS2A16XpeB4paDyjsltVHv+6azoA6wAY8WLORYd6jGMAq4zgcTd/6BArojm/prhViwo8tsQKFRSAq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1Cq0gE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dyswCEWE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dwlAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAACEWKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQolAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAAAEXIFCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAARggR4Li5f/hJviWsPse5R7SzU/wp7r8u4szV3KnW5FahpAAAA=="
	TEST_54155_PSBT_FINAL_TXHEX    = "0200000000010183fd0482b2c763cc5a116e3004f7c1c3834af5916a43c8d56bc7550126c007650000000000fdffffff01282300000000000016001450dceca158a9c872eb405d52293d351110572c9e04401e6a744fd9c20321247b9450bfe8ac82acdd008d0671646002a1aff29545601f5b2640da8f212a5192b14ed8a024af86bb4038f0fdeff33d689b56deb7b33bef40f49fb85f536a0b3724b435bcad6c98a0a087060671647800000b87f99beb234527711da9a526237a565bddf12da5fe1ea53250428d286d03339f0563a680cb1244202ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350aad201387aab21303782b17e760c670432559df3968e52cb82cc2d8f9be43a227d5dcac41c150929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0063c58b39161dea318c02ae3381c4ddffa040ae88e6fe9ae1562c28f2db1028500000000"
)

func TestMain(m *testing.M) {

	cfg, err := config.GetConfig("../../../example/config-test.yaml")
	if err != nil {
		panic(err)
	}

	parsedConfig, err := cfg.Parse()
	if err != nil {
		panic(err)
	}

	var broadcaster btc.BtcClientInterface

	// if cfg.BtcNodeConfig.Network == "testnet4" {
	// 	fmt.Println("Using raw rpc client for testnet4")
	// 	broadcaster, err = btc.NewRawRpcClient(cfg.BtcNodeConfig.Host, cfg.BtcNodeConfig.User, cfg.BtcNodeConfig.Pass, cfg.BtcNodeConfig.Network)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else {
	// 	broadcaster, err = btc.NewBtcClient(parsedConfig.BtcNodeConfig)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	broadcaster, err = btc.NewBtcClient(parsedConfig.BtcNodeConfig, cfg.BtcNodeConfig.Network)
	if err != nil {
		panic(err)
	}

	signerClient, err := btc.NewBtcClient(parsedConfig.BtcSignerConfig, cfg.BtcSignerConfig.Network)
	if err != nil {
		panic(err)
	}

	signer := btc.NewPsbtSigner(signerClient, parsedConfig.BtcSignerConfig.Address, parsedConfig.BtcSignerConfig.Passphrase, parsedConfig.BtcSignerConfig.Network)

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
	//userSignedPsbtBase64 := "cHNidP8BAFICAAAAAfUimJvNwyq+s+JWNEHbilDru0GNci6pMmBiGZvpGT74AAAAAAD9////ASgjAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKxAnAAAAAAAAIlEg2t54XUPHU7zIxm8h/vBWQ+u02YEqpgeCx0QBiSVbu0sBAwQAAAAAQRQq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1CoshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfQPNNqOPdOT+My/MPUElUddDtW+RLSjTsOZy3BAPkSZpb6iaYnl3kpv2uXYNSkic4U2x4uTjMZYdkw9xoa+NqCehCFcBQkpt0waBJVLeLS2A16XpeB4paDyjsltVHv+6azoA6wLA+TxG6WUpeNIqF9MLRbzubGb4+7/SUwSwwUBU1hSVQRSAq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1Cq0gE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dyswCEWE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dwlAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAACEWKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQolAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAAAEXIFCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAARggFq1jps0IRdYmx2Bk96iRBq0h1JkJ5fGwYFdUZNhvlgYAAA=="
	userSignedPsbtBase64 := UNSTAKE_54167_PSBT_B64
	finalTxHex := UNSTAKE_54167_PSBT_FINAL_TXHEX
	packet, err := psbt.NewFromRawBytes(strings.NewReader(userSignedPsbtBase64), true)
	if err != nil {
		t.Fatalf("Failed to parse PSBT: %v", err)
	}

	finalTx, err := mockHandler.SignPsbt(packet)
	if err != nil {
		t.Fatalf("Failed to sign PSBT: %v", err)
	}

	buf := new(bytes.Buffer)
	finalTx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())

	t.Logf("finalTx: %s\n", txHex)

	if txHex != finalTxHex {
		t.Fatalf("Final txHex is not match")
	}

	testResults, err := mockHandler.TestMempoolAccept([]*wire.MsgTx{finalTx}, 0.1)
	if err != nil {
		t.Fatalf("Failed to test mempool accept: %v", err)
	}
	assert.NoError(t, err)
	t.Logf("TestMempoolAcceptResults: %+v\n", testResults[0])

	txid, err := mockHandler.BroadcastTx(finalTx)
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("Failed to send raw transaction: %v", err)
	} else {
		t.Logf("txid: %s\n", txid)
	}
}
