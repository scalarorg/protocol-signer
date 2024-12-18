package cmd

import (
	"github.com/spf13/cobra"

	"github.com/scalarorg/protocol-signer/config"
	signerservice "github.com/scalarorg/protocol-signer/internals/signer"
	"github.com/scalarorg/protocol-signer/packages/btc"
)

func init() {
	rootCmd.AddCommand(runSignerCmd)
}

var runSignerCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the signer service",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString(configPathKey)
		if err != nil {
			return err
		}

		cfg, err := config.GetConfig(configPath)
		if err != nil {
			return err
		}

		parsedConfig, err := cfg.Parse()
		if err != nil {
			return err
		}

		// var broadcaster btc.BtcClientInterface

		// if cfg.BtcNodeConfig.Network == "testnet4" {
		// 	fmt.Println("Using raw rpc client for testnet4")
		// 	broadcaster, err = btc.NewRawRpcClient(cfg.BtcNodeConfig.Host, cfg.BtcNodeConfig.User, cfg.BtcNodeConfig.Pass, cfg.BtcNodeConfig.Network)
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {
		// 	broadcaster, err = btc.NewBtcClient(parsedConfig.BtcNodeConfig)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		broadcaster, err := btc.NewBtcClient(parsedConfig.BtcNodeConfig, cfg.BtcNodeConfig.Network)
		if err != nil {
			return err
		}

		signerClient, err := btc.NewBtcClient(parsedConfig.BtcSignerConfig, cfg.BtcSignerConfig.Network)
		if err != nil {
			return err
		}
		// TODO: Add options to use customn remote signers
		// Integrate cubist remote signer
		signer := btc.NewPsbtSigner(signerClient, parsedConfig.BtcSignerConfig.Address, parsedConfig.BtcSignerConfig.Passphrase, parsedConfig.BtcSignerConfig.Network)

		srv, err := signerservice.New(
			cmd.Context(),
			parsedConfig,
			signer,
			broadcaster,
		)

		if err != nil {
			return err
		}

		// TODO: Add signal handling and gracefull shutdown
		return srv.Start()
	},
}
