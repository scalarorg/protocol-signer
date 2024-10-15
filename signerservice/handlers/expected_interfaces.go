package handlers

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type ExternalEvmClient interface {
	ChainName() string
	CheckUnbondingTx(ctx context.Context, txHash common.Hash, unbondingTx *string) error
}
