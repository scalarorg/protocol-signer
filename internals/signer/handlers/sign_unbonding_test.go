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
	userSignedPsbtBase64 := "cHNidP8BAFICAAAAAfUimJvNwyq+s+JWNEHbilDru0GNci6pMmBiGZvpGT74AAAAAAD9////ASgjAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKxAnAAAAAAAAIlEg2t54XUPHU7zIxm8h/vBWQ+u02YEqpgeCx0QBiSVbu0sBAwQAAAAAQRQq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1CoshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfQPNNqOPdOT+My/MPUElUddDtW+RLSjTsOZy3BAPkSZpb6iaYnl3kpv2uXYNSkic4U2x4uTjMZYdkw9xoa+NqCehCFcBQkpt0waBJVLeLS2A16XpeB4paDyjsltVHv+6azoA6wLA+TxG6WUpeNIqF9MLRbzubGb4+7/SUwSwwUBU1hSVQRSAq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1Cq0gE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dyswCEWE4eqshMDeCsX52DGcEMlWd85aOUsuCzC2Pm+Q6In1dwlAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAACEWKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQolAYshIJihyflfrfabq/5zjDSJchXpFwfx/bqZ+lR02TsfAAAAAAEXIFCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAARggFq1jps0IRdYmx2Bk96iRBq0h1JkJ5fGwYFdUZNhvlgYAAA=="

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

	if txHex != "02000000000101f522989bcdc32abeb3e2563441db8a50ebbb418d722ea9326062199be9193ef80000000000fdffffff01282300000000000016001450dceca158a9c872eb405d52293d351110572c9e0440bf5bd22e974a3c7b4e40c93343a39f69fd417f9fbb8df7901d12fe3d11b57b09dced0fdf19f78e267d82f0edc6be790270a2abc3df9f082ec1ced4c40828b91240f34da8e3dd393f8ccbf30f50495475d0ed5be44b4a34ec399cb70403e4499a5bea26989e5de4a6fdae5d8352922738536c78b938cc658764c3dc686be36a09e844202ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350aad201387aab21303782b17e760c670432559df3968e52cb82cc2d8f9be43a227d5dcac41c050929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0b03e4f11ba594a5e348a85f4c2d16f3b9b19be3eeff494c12c3050153585255000000000" {
		t.Fatalf("txHex is empty")
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

func TestUnstakingSignUserProtocol(t *testing.T) {
	//https://mempool.space/testnet4/api/tx/10f5d2f7167428cfd983bfbaad566adce246f98d3a0ca8ab590844bcab9b2c81/hex
	txHex := "02000000000101357229066a1315a3a411517bf8db9ae31745e14f6c66ac4c807e5860baa334960200000000fdffffff035504000000000000225120dade785d43c753bcc8c66f21fef05643ebb4d9812aa60782c7440189255bbb4b00000000000000003d6a013504531801040100080000000000aa36a714b91e3a8ef862567026d6f376c9f3d6b814ca43371424a1db57fa3ecafcbad91d6ef068439aceeae090d60a8f000000000016001450dceca158a9c872eb405d52293d351110572c9e02483045022100be4ec4a872b521f4415ad8463a90d28258c32eba741fa8b792260863e5dba11502203e3c67cac9b33a1f705fd4ef5215530c85bc1817306af515b0870ca94c97c3500121022ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350a00000000"
	packet, err := psbt.NewFromRawBytes(strings.NewReader(txHex), false)
	if err != nil {
		t.Fatalf("Failed to parse PSBT: %v", err)
	}
	t.Logf("packet: %+v\n", packet)

	// finalTx, err := mockHandler.SignPsbt(packet)
	// assert.NoError(t, err)

	// buf := new(bytes.Buffer)
	// finalTx.Serialize(buf)
	// txHex = hex.EncodeToString(buf.Bytes())
	// t.Logf("txHex: %s\n", txHex)
}
