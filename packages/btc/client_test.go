package btc_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/scalarorg/protocol-signer/config"
	"github.com/scalarorg/protocol-signer/internals/signer/handlers"
	"github.com/scalarorg/protocol-signer/packages/btc"
	"github.com/stretchr/testify/assert"
)

var mockHandler *handlers.Handler
var broadcaster btc.BtcClientInterface

const (
	STAKE_FINAL_PSBT_HEX = "02000000000101f1c479e7820eb5c7f85acc311f05b17add5f3040df63fa619ef5087dfcc08f8b0200000000fdffffff03102700000000000022512067bff357780a93826a444646aec681c4ff1f4316244478c0d611f91a75c93b8a00000000000000003d6a013504531801040100080000000000aa36a714b91e3a8ef862567026d6f376c9f3d6b814ca43371424a1db57fa3ecafcbad91d6ef068439aceeae090020280000000000016001450dceca158a9c872eb405d52293d351110572c9e024830450221009c5b34dc30a22059e6415c0ff4c819f2b091f2deb40377f7b0f603b8c7ad12ac0220483c73879c5ac3293c003b42835834f634546fa26a6e6809c2c8df993e325dec0121022ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350a00000000"
)

func TestMain(m *testing.M) {

	cfg, err := config.GetConfig("../../example/config-test.yaml")
	if err != nil {
		panic(err)
	}

	parsedConfig, err := cfg.Parse()
	if err != nil {
		panic(err)
	}

	broadcaster, err = btc.NewBtcClient(parsedConfig.BtcNodeConfig, cfg.BtcNodeConfig.Network)
	if err != nil {
		panic(err)
	}

	m.Run()
}

// Note: To run this test, must build bitcoin-vault-ffi first then copy to the lib folder
// cp -p ../../bitcoin-vault/target/release/libbitcoin_vault_ffi.* ./lib

// CGO_LDFLAGS="-L./lib -lbitcoin_vault_ffi" CGO_CFLAGS="-I./lib" go test -timeout 10m -run ^TestSignUnBondingIngoreCheckOnEvm$ github.com/scalarorg/protocol-signer/internals/signer/handlers -v -count=1
func TestSignUnBondingIngoreCheckOnEvm(t *testing.T) {
	finalTxHex, err := hex.DecodeString(STAKE_FINAL_PSBT_HEX)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	var finalTx wire.MsgTx
	err = finalTx.Deserialize(bytes.NewReader(finalTxHex))
	if err != nil {
		t.Fatalf("Failed to parse tx: %v", err)
	}
	serializedTx, err := btc.CreateRawTx(&finalTx)
	if err != nil {
		t.Fatalf("Failed to create raw tx: %v", err)
	}
	assert.Equal(t, serializedTx, STAKE_FINAL_PSBT_HEX)
	testResults, err := broadcaster.TestMempoolAccept([]*wire.MsgTx{&finalTx}, 0.1)
	if err != nil {
		t.Fatalf("Failed to test mempool accept: %v", err)
	}
	assert.NoError(t, err)
	t.Logf("TestMempoolAcceptResults: %+v\n", testResults[0])

	txid, err := broadcaster.SendTx(&finalTx)
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("Failed to send raw transaction: %v", err)
	} else {
		t.Logf("txid: %s\n", txid)
	}
}
