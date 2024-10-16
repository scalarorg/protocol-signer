package btc_test

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"log"
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

	t.Logf(">>> Before: %+v\n", packet.Inputs)

	result, err := btc.SignPsbtAll(packet, privKey)
	if err != nil {
		t.Fatalf("Failed to sign PSBT: %v", err)
	}

	t.Logf(">>> After: %+v\n", packet.Inputs)

	if result == nil {
		t.Fatalf("No result returned from signing")
	}

	err = psbt.MaybeFinalizeAll(packet)
	if err != nil {
		t.Fatalf("Failed to finalize PSBT: %v", err)
	}

	t.Logf(">>> Finalized: %+v\n", packet.Inputs)

	finalTx, err := psbt.Extract(packet)
	if err != nil {
		t.Fatalf("Failed to extract transaction: %v", err)
	}

	t.Logf(">>> Extracted: %v\n", finalTx)

	txHash := finalTx.TxHash()
	t.Logf(">>> TxHash: %v\n", txHash)

	// Serialize the transaction to raw bytes
	var buf bytes.Buffer
	finalTx.Serialize(&buf)

	// Convert to hex
	hexTx := hex.EncodeToString(buf.Bytes())
	log.Printf(">>> Transaction Hex: %s\n", hexTx)

	// Convert to base64
	base64Tx := base64.StdEncoding.EncodeToString(buf.Bytes())
	log.Printf(">>> Transaction Base64: %s\n", base64Tx)
}

// 4 64 112 183 15 112 5 51 122 103 2 190 156 196 163 97 80 43 83 67 64 158 252 252 65 234 34 85 238 172 83 65 171 135 218 187 172 216 110 25 215 159 3 232 113 166 221 121 170 107 249 103 193 197 110 15 54 200 151 200 141 165 35 135 213 207 64 243 214 87 225 251 105 114 184 153 45 249 96 85 210 114 150 254 65 161 133 211 68 134 183 122 230 179 252 216 124 218 49 133 91 122 172 67 11 21 61 189 180 48 55 145 31 204 246 240 16 222 59 24 131 65 220 90 166 170 94 17 178 22 2 190

// 68 32 42 227 30 168 112 154 237 168 25 75 163 226 247 231 233 94 104 14 139 101 19 92 137 131 192 162 152 209 123 197 53 10 173 32 207 93 255 87 161 115 197 172 131 35 196 186 202 63 255 7 40 235 113 111 57 240 229 166 3 18 50 12 210 147 91 12 172 97 193 80 146 155 116 193 160 73 84 183 139 75 96 53 233 122 94 7 138 90 15 40 236 150 213 71 191 238 154 206 128 58 192 120 128 80 231 157 83 6 55 178 191 150 62 199 158 115 158 164 120 151 139 119 179 98 100 148 57 226 0 69 205 203 86 110 229 51 71 188 235 230 196 197 47 11 25 75 138 195 165 143 235 224 209 172 101 34 124 123 75 20 32 238 73 17 204

// 4 64 243 214 87 225 251 105 114 184 153 45 249 96 85 210 114 150 254 65 161 133 211 68 134 183 122 230 179 252 216 124 218 49 133 91 122 172 67 11 21 61 189 180 48 55 145 31 204 246 240 16 222 59 24 131 65 220 90 166 170 94 17 178 222 190 64 197 188 244 161 140 202 134 81 210 76 232 173 112 198 180 166 13 202 135 74 12 174 1 180 22 87 105 85 106 65 29 255 125 173 139 183 230 55 202 218 169 190 20 119 183 182 158 121 155 209 247 183 170 184 90 221 61 91 110 105 190 34 171 84

// 68 32 42 227 30 168 112 154 237 168 25 75 163 226 247 231 233 94 104 14 139 101 19 92 137 131 192 162 152 209 123 197 53 10 173 32 207 93 255 87 161 115 197 172 131 35 196 186 202 63 255 7 40 235 113 111 57 240 229 166 3 18 50 12 210 147 91 12 172 97 193 80 146 155 116 193 160 73 84 183 139 75 96 53 233 122 94 7 138 90 15 40 236 150 213 71 191 238 154 206 128 58 192 120 128 80 231 157 83 6 55 178 191 150 62 199 158 115 158 164 120 151 139 119 179 98 100 148 57 226 0 69 205 203 86 110 229 51 71 188 235 230 196 197 47 11 25 75 138 195 165 143 235 224 209 172 101 34 124 123 75 20 32 238 73 17 204
