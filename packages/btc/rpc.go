package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

type RawRpcClient struct {
	host    string
	user    string
	pass    string
	network string
}

func NewRawRpcClient(host, user, pass, network string) (*RawRpcClient, error) {
	return &RawRpcClient{host: host, user: user, pass: pass, network: network}, nil
}

type RPCResponse struct {
	Result interface{} `json:"result"`
	Error  *RPCError   `json:"error"`
	ID     string      `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *RawRpcClient) SendTx(tx *wire.MsgTx) (*chainhash.Hash, error) {

	buf := bytes.NewBuffer(nil)
	tx.Serialize(buf)

	txHex := hex.EncodeToString(buf.Bytes())

	allowHighFees := true
	params := []interface{}{txHex, allowHighFees}
	response := &RPCResponse{}

	err := c.sendRequest("sendrawtransaction", params, response)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %+v\n", response)

	if response.Error != nil {
		return nil, fmt.Errorf("error sending tx: %s", response.Error.Message)
	}

	txidStr, ok := response.Result.(string)
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", response.Result)
	}

	return chainhash.NewHashFromStr(txidStr)
}

func (c *RawRpcClient) sendRequest(method string, params []interface{}, response interface{}) error {
	payload := map[string]interface{}{
		"jsonrpc": "1.0",
		"id":      "curltest",
		"method":  method,
		"params":  params,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s", c.host)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	fmt.Printf("Request: %s\n", req.URL)

	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Raw response: %s\n", string(bodyBytes))

	return json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(response)
}
