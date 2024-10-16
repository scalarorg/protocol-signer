package evm

import (
	"context"
	"errors"
	"fmt"

	"github.com/axelarnetwork/utils/monads/results"
	"github.com/axelarnetwork/utils/slices"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

func (e *EvmClient) CheckUnbondingTx(ctx context.Context, txHash common.Hash, psbtBase64 string) error {

	response, err := e.TransactionReceipts(context.Background(), []common.Hash{txHash})
	if err != nil {
		return fmt.Errorf("cannot get transaction receipts: %w", err)
	}

	if len(response) != 1 {
		return fmt.Errorf("expected 1 transaction receipt, got %d", len(response))
	}

	receipt := results.Result[ethTypes.Receipt](response[0]).Ok()

	eventSignatureHash := crypto.Keccak256Hash(ContractCallEventSignature)

	for _, log := range receipt.Logs {
		if log.Topics[0] == eventSignatureHash {
			event, err := DecodeEventContractCall(log)
			if err != nil {
				return fmt.Errorf("failed to decode event: %w", err)
			}

			decodedPsbt, err := DecodePsbt(event.Payload)
			if err != nil {
				return fmt.Errorf("failed to decode payload: %w", err)
			}

			if decodedPsbt != psbtBase64 {
				return errors.New("psbt does not match")
			}

		}
	}

	return nil
}

func DecodeEventContractCall(log *ethTypes.Log) (*EventContractCall, error) {
	// Create an event instance to store unpacked values
	event := &EventContractCall{}

	// Check and decode the indexed parameters
	if len(log.Topics) < 3 {
		return nil, fmt.Errorf("log does not contain enough topics for event")
	}

	// Unpack indexed parameters
	event.Sender = common.HexToAddress(log.Topics[1].Hex())
	payloadHash := common.BytesToHash(log.Topics[2].Bytes()) // Adjusted to bytes

	// Set the PayloadHash in the event
	copy(event.PayloadHash[:], payloadHash[:]) // Copy the hash into the array

	// Now unpack the non-indexed parameters
	err := parsedABI.UnpackIntoInterface(event, "ContractCall", log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack event data: %w", err)
	}

	if len(event.Payload) < 64 {
		return nil, fmt.Errorf("invalid encoded data")
	}

	return event, nil
}

func DecodePsbt(payload []byte) (string, error) {
	if len(payload) < 32 {
		return "", errors.New("payload is too short to contain valid ABI-encoded data")
	}

	// Create a new ABI Type for a single string parameter
	psbtType, err := abi.NewType("string", "", nil)
	if err != nil {
		return "", err
	}

	// Unpack the ABI-encoded string
	decoded, err := abi.Arguments{{Type: psbtType}}.Unpack(payload)
	if err != nil {
		return "", err
	}

	if len(decoded) == 0 {
		return "", errors.New("no values unpacked")
	}

	decodedString, ok := decoded[0].(string)
	if !ok {
		return "", errors.New("decoded value is not a string")
	}

	return decodedString, nil
}

func (c *EvmClient) TransactionReceipts(ctx context.Context, txHashes []common.Hash) ([]TxReceiptResult, error) {
	batch := slices.Map(txHashes, func(txHash common.Hash) rpc.BatchElem {
		var receipt *ethTypes.Receipt
		return rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{txHash},
			Result: &receipt,
		}
	})

	if err := c.rpc.BatchCallContext(ctx, batch); err != nil {
		return nil, fmt.Errorf("unable to send batch request: %v", err)
	}

	return slices.Map(batch, func(elem rpc.BatchElem) TxReceiptResult {
		if elem.Error != nil {
			return TxReceiptResult(results.FromErr[ethTypes.Receipt](elem.Error))
		}

		receipt := elem.Result.(**ethTypes.Receipt)
		if *receipt == nil {
			return TxReceiptResult(results.FromErr[ethTypes.Receipt](ethereum.NotFound))
		}

		return TxReceiptResult(results.FromOk(**receipt))
	}), nil

}
