package evm

type EvmConfig struct {
	ChainId              string           `mapstructure:"chain-id"`
	ChainName            string           `mapstructure:"chain-name"`
	RpcUrl               string           `mapstructure:"rpc-url"`
	FinalityOverride     FinalityOverride `mapstructure:"finality-override"`
	SmartContractAddress string           `mapstructure:"smart-contract-address"`
}

func DefaultEvmConfig() *EvmConfig {
	return &EvmConfig{
		ChainId:              "1337",
		ChainName:            "evm-local",
		RpcUrl:               "http://localhost:8545",
		FinalityOverride:     NoOverride,
		SmartContractAddress: "",
	}
}
