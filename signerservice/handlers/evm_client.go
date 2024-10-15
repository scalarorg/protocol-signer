package handlers

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/protocol-signer/config"
	"github.com/scalarorg/protocol-signer/evmclient"
)

type EvmClient struct {
	evmConfig config.EvmConfig
	client    evmclient.Client
}

var _ ExternalEvmClient = (*EvmClient)(nil)

func NewEvmClient(evmConfig config.EvmConfig) ExternalEvmClient {
	return &EvmClient{evmConfig: evmConfig}
}

func (e *EvmClient) ChainName() string {
	return e.evmConfig.ChainName
}

func (e *EvmClient) CheckUnbondingTx(ctx context.Context, txHash common.Hash, unbondingTx *string) error {
	client, err := e.GetClient()
	if err != nil {
		return err
	}
	receipts, err := client.TransactionReceipts(ctx, []common.Hash{txHash})
	if err != nil {
		return err
	}

	if len(receipts) == 0 {
		return errors.New("transaction not found")
	}
	// receipt := receipts[0]
	// Todo: Check if the receipt is matching the unbondingTx
	// Parse unbondingTx to get EvmTxId
	return nil
}

func (e *EvmClient) GetClient() (evmclient.Client, error) {
	if e.client == nil {
		client, err := evmclient.NewClient(e.evmConfig.RpcUrl, e.evmConfig.FinalityOverride)
		if err != nil {
			return nil, err
		}
		e.client = client
	}
	return e.client, nil
}
