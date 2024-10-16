package evm_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/protocol-signer/packages/evm"
	"github.com/stretchr/testify/require"
)

var mockConfig evm.EvmConfig

func TestMain(m *testing.M) {
	mockConfig = evm.EvmConfig{
		RpcUrl:           "https://eth-sepolia.g.alchemy.com/v2/nNbspp-yjKP9GtAcdKi8xcLnBTptR2Zx",
		ChainName:        "ethereum-sepolia",
		FinalityOverride: evm.Confirmation,
	}
	m.Run()
}

func TestDecodePayload(t *testing.T) {

	var psbtBase64 = "cHNidP8BAFICAAAAAVrQluj7zrZI/RJ4h3pXKY5UQckb3hm1I7h+rwqUVXwJAAAAAAD9////AXImAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKw8nAAAAAAAAIlEgkFj0s5tPU4QO+QWJFepTzvTiZxWVJMmxkN+zd/AYnZsBCP0qAQRAcLcPcAUzemcCvpzEo2FQK1NDQJ78/EHqIlXurFNBq4fau6zYbhnXnwPocabdeapr+WfBxW4PNsiXyI2lI4fVz0Dz1lfh+2lyuJkt+WBV0nKW/kGhhdNEhrd65rP82HzaMYVbeqxDCxU9vbQwN5EfzPbwEN47GINB3Fqmql4Rst6+RCAq4x6ocJrtqBlLo+L35+leaA6LZRNciYPAopjRe8U1Cq0gz13/V6FzxayDI8S6yj//ByjrcW858OWmAxIyDNKTWwysYcFQkpt0waBJVLeLS2A16XpeB4paDyjsltVHv+6azoA6wHiAUOedUwY3sr+WPseec56keJeLd7NiZJQ54gBFzctWbuUzR7zr5sTFLwsZS4rDpY/r4NGsZSJ8e0sUIO5JEcwAAA=="

	var psbtHex = "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000025063484e69645038424146494341414141415672516c756a377a725a492f524a34683370584b59355551636b6233686d314937682b727771555658774a41414141414144392f2f2f2f4158496d414141414141414146674155554e7a736f56697079484c72514631534b543031455242584c4a344141414141414145424b77386e4141414141414141496c45676b466a30733574505534514f2b51574a466570547a7654695a7857564a4d6d786b4e2b7a642f41596e5a73424350307141515241634c63506341557a656d634376707a456f3246514b314e44514a37382f454871496c587572464e427134666175367a5962686e586e77506f63616264656170722b576642785734504e7369587949326c493466567a30447a316c66682b326c79754a6b742b574256306e4b572f6b476868644e45687264363572503832487a614d5956626571784443785539766251774e3545667a506277454e343747494e423346716d716c34527374362b524341713478366f634a727471426c4c6f2b4c33352b6c656141364c5a524e63695950416f706a5265385531437130677a31332f5636467a7861794449385336796a2f2f42796a7263573835384f576d41784979444e4b5457777973596346516b7074307761424a564c654c53324131365870654234706144796a736c745648762b36617a6f413677486941554f65645577593373722b57507365656335366b654a654c64374e695a4a5135346742467a6374576275557a52377a72357354464c77735a5334724470592f72344e47735a534a3865307355494f354a4563774141413d3d00000000000000000000000000000000"

	payloadToBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	decodedPayload, err := evm.DecodePsbt(payloadToBytes)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if decodedPayload != psbtBase64 {
		t.Errorf("Decoded Payload does not match")
	}
}

func TestDecodePsbtBase64(t *testing.T) {
	stakerPsbt := "cHNidP8BAFICAAAAAVrQluj7zrZI/RJ4h3pXKY5UQckb3hm1I7h+rwqUVXwJAAAAAAD9////AXImAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBKw8nAAAAAAAAIlEgkFj0s5tPU4QO+QWJFepTzvTiZxWVJMmxkN+zd/AYnZtBFCrjHqhwmu2oGUuj4vfn6V5oDotlE1yJg8CimNF7xTUKa/PeGFPFI+tlh34WOxsYSsPX8wg4uJXl35HhYElEst9A89ZX4ftpcriZLflgVdJylv5BoYXTRIa3euaz/Nh82jGFW3qsQwsVPb20MDeRH8z28BDeOxiDQdxapqpeEbLevmIVwVCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAeIBQ551TBjeyv5Y+x55znqR4l4t3s2JklDniAEXNy1Zu5TNHvOvmxMUvCxlLisOlj+vg0axlInx7SxQg7kkRzEUgKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQqtIM9d/1ehc8WsgyPEuso//wco63FvOfDlpgMSMgzSk1sMrMAAAA=="

	testPsbt, err := psbt.NewFromRawBytes(strings.NewReader(stakerPsbt), true)
	if err != nil {
		t.Fatalf("Unable to parse Psbt: %v", err)
	}
	t.Logf("Successfully parsed, got transaction: %v", spew.Sdump(testPsbt))

	var b bytes.Buffer
	err = testPsbt.Serialize(&b)
	if err != nil {
		t.Fatalf("Unable to serialize created Psbt: %v", err)
	}

	base64Packet := base64.StdEncoding.EncodeToString(b.Bytes())
	require.Equal(t, stakerPsbt, base64Packet, "Base64 encoded Psbt does not match")
}

func TestCheckUnbondingTx(t *testing.T) {
	mockClient := evm.NewEvmClient(mockConfig)
	var txID = "0xb7ced6fe1b05b01a9f24ae84a14890ea2fd6ef1ff42ae82fe5110ff649e9dc41"
	var psbtBase64 = "cHNidP8BAFICAAAAAUM2R8LOv9D1acCrQC4g9reCh+3u0CkprhQ/CGrzYFrKAAAAAAD9////ASciAAAAAAAAFgAUUNzsoVipyHLrQF1SKT01ERBXLJ4AAAAAAAEBK8QiAAAAAAAAIlEgUiPouA919Arm/mld5UVulY4HFUe3CUvidmtDCJyUOHRBFCrjHqhwmu2oGUuj4vfn6V5oDotlE1yJg8CimNF7xTUKGqpwf1oaV/8WEF+uunL1V6iWOpitlLONjh+Oud3IvilAKa0/VNbnkDtjEAOElCNwEcnuJCYlGdwjfPx+GIgzYU8VFQeXZL+Fz8NvOClV73QpNcdjOZGnAHhJ+Mm9sltC52IVwVCSm3TBoElUt4tLYDXpel4HiloPKOyW1Ue/7prOgDrAFuIuioNMBacc5AXCkPpmmlMAQJR7esmfH/PGLWmWX25mR0dHHl1WKnImL8iqD3X/VURzLq83wrCHDtY3AZcpEkUgKuMeqHCa7agZS6Pi9+fpXmgOi2UTXImDwKKY0XvFNQqtIPNYZ068zmiDKsMCuC7Hzk6c3Y+FZoTKPVRWhqJAvhXYrMAAAA=="

	ctx := context.Background()
	txHash := common.HexToHash(txID)
	err := mockClient.CheckUnbondingTx(ctx, txHash, psbtBase64)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
