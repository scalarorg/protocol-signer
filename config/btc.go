package config

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
)

type BtcConfig struct {
	Host       string `mapstructure:"host"`
	User       string `mapstructure:"user"`
	Pass       string `mapstructure:"pass"`
	Network    string `mapstructure:"network"`
	Address    string `mapstructure:"address"`
	Passphrase string `mapstructure:"passphrase"`
}

type ParsedBtcConfig struct {
	Host       string
	User       string
	Pass       string
	Network    *chaincfg.Params
	Address    string
	Passphrase string
}

func DefaultBtcConfig() *BtcConfig {
	return &BtcConfig{
		Host:       "localhost:18556",
		User:       "user",
		Pass:       "pass",
		Network:    "regtest",
		Address:    "",
		Passphrase: "",
	}
}

func (c *BtcConfig) Parse() (*ParsedBtcConfig, error) {
	params, err := c.getBtcNetworkParams()

	if err != nil {
		return nil, err
	}
	return &ParsedBtcConfig{
		Host:       c.Host,
		User:       c.User,
		Pass:       c.Pass,
		Network:    params,
		Address:    c.Address,
		Passphrase: c.Passphrase,
	}, nil
}

func (cfg *BtcConfig) getBtcNetworkParams() (*chaincfg.Params, error) {
	switch cfg.Network {
	case "testnet4":
		return &chaincfg.TestNet3Params, nil // TODO: Change to TestNet4Params
	case "testnet3":
		return &chaincfg.TestNet3Params, nil
	case "mainnet":
		return &chaincfg.MainNetParams, nil
	case "regtest":
		return &chaincfg.RegressionNetParams, nil
	case "simnet":
		return &chaincfg.SimNetParams, nil
	case "signet":
		return &chaincfg.SigNetParams, nil
	default:
		return nil, fmt.Errorf("unknown network %s", cfg.Network)
	}
}
