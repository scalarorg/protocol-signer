package evm

import (
	"fmt"
	"strings"

	"github.com/axelarnetwork/utils/monads/results"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxReceiptResult results.Result[types.Receipt]

type FinalityOverride int

const (
	NoOverride FinalityOverride = iota
	Confirmation
)

func ParseFinalityOverride(s string) (FinalityOverride, error) {
	switch strings.ToLower(s) {
	case "":
		return NoOverride, nil
	case strings.ToLower(string(Confirmation.String())):
		return Confirmation, nil
	default:
		return -1, fmt.Errorf("invalid finality override option")
	}
}

// String returns the string representation of the FinalityOverride
func (fo FinalityOverride) String() string {
	switch fo {
	case NoOverride:
		return "NoOverride"
	case Confirmation:
		return "Confirmation"
	default:
		return "Unknown"
	}
}

type EventContractCall struct {
	Sender                     common.Address `json:"sender"`
	DestinationChain           string         `json:"destinationChain"`
	DestinationContractAddress string         `json:"destinationContractAddress"`
	PayloadHash                [32]byte       `json:"payloadHash"`
	Payload                    []byte         `json:"payload"`
}
