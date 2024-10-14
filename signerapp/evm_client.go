package signerapp

import (
	"github.com/scalarorg/protocol-signer/config"
	"github.com/scalarorg/protocol-signer/evmclient"
)

type EvmClient struct {
	evmConfig config.EvmConfig
	client    *evmclient.Client
}

var _ ExternalEvmClient = (*EvmClient)(nil)

func NewEvmClient(evmConfig config.EvmConfig) ExternalEvmClient {
	return &EvmClient{evmConfig: evmConfig}
}

func (e *EvmClient) GetClient() (*evmclient.Client, error) {
	if e.client == nil {
		client, err := evmclient.NewClient(e.evmConfig.RpcUrl, e.evmConfig.FinalityOverride)
		if err != nil {
			return nil, err
		}
		e.client = &client
	}
	return e.client, nil
}
