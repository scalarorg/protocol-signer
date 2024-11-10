package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

func CreateRawTx(tx *wire.MsgTx) (string, error) {
	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	// TODO(yy): add similar checks found in `BtcDecode` to
	// `BtcEncode` - atm it just serializes bytes without any
	// bitcoin-specific checks.
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	// Sanity check the provided tx is valid, which can be removed
	// once we have similar checks added in `BtcEncode`.
	//
	// NOTE: must be performed after buf.Bytes is copied above.
	//
	// TODO(yy): remove it once the above TODO is addressed.
	if err := tx.Deserialize(buf); err != nil {
		err = fmt.Errorf("%w: %v", rpcclient.ErrInvalidParam, err)
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func CreateRawTxs(txns []*wire.MsgTx) ([]string, error) {
	// Iterate all the transactions and turn them into hex strings.
	rawTxns := make([]string, 0, len(txns))
	for _, tx := range txns {
		rawTx, err := CreateRawTx(tx)
		if err != nil {
			return nil, err
		}
		rawTxns = append(rawTxns, rawTx)

	}

	return rawTxns, nil
}
