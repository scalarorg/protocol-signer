package evm

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var ContractCallEventSignature = []byte("ContractCall(address,string,string,bytes32,bytes)")

var ContractCallEventABI = `[{
	"anonymous":false,
	"inputs":[
		{"indexed":true,"name":"sender","type":"address"},
		{"indexed":false,"name":"destinationChain","type":"string"},
		{"indexed":false,"name":"destinationContractAddress","type":"string"},
		{"indexed":true,"name":"payloadHash","type":"bytes32"},
		{"indexed":false,"name":"payload","type":"bytes"}
	],
	"name":"ContractCall",
	"type":"event"
}]`

// Get the ABI definition of the contract
var parsedABI, _ = abi.JSON(strings.NewReader(ContractCallEventABI))

type EvmClient struct {
	evmConfig EvmConfig
	raw       *ethclient.Client
	rpc       *rpc.Client
}

func NewEvmClient(config EvmConfig) *EvmClient {
	rpc, err := rpc.DialContext(context.Background(), config.RpcUrl)
	if err != nil {
		panic("failed to create an RPC connection for EVM chain. Verify your RPC config.")
	}

	_, err = ethclient.NewClient(rpc).BlockNumber(context.Background())
	if err != nil {
		panic("failed to create an EVM client. Verify your RPC config.")
	}

	return &EvmClient{
		evmConfig: config,
		raw:       ethclient.NewClient(rpc),
		rpc:       rpc,
	}
}

func (e *EvmClient) ChainName() string {
	return e.evmConfig.ChainName
}
