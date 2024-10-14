package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scalarorg/protocol-signer/btcclient"
	"github.com/scalarorg/protocol-signer/config"
	m "github.com/scalarorg/protocol-signer/observability/metrics"
	"github.com/scalarorg/protocol-signer/signerapp"
	"github.com/scalarorg/protocol-signer/signerservice"
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

		parsedGlobalParams, err := signerapp.NewVersionedParamsRetriever(globalParamPath)

		if err != nil {
			return err
		}

		fullNodeClient, err := btcclient.NewBtcClient(parsedConfig.BtcNodeConfig)

		if err != nil {
			return err
		}

		chainInfo := signerapp.NewBitcoindChainInfo(fullNodeClient)

		signerClient, err := btcclient.NewBtcClient(parsedConfig.BtcSignerConfig)
		// evmClient, err := evmclient.NewEvmClient(parsedConfig.EvmConfigs[0])
		if err != nil {
			return err
		}
		// TODO: Add options to use customn remote signers
		// Integrate cubist remote signer
		signer := signerapp.NewPsbtSigner(signerClient)

		app := signerapp.NewSignerApp(
			signer,
			chainInfo,
			parsedGlobalParams,
			parsedConfig.BtcNodeConfig.Network,
		)

		metrics := m.NewCovenantSignerMetrics()

		srv, err := signerservice.New(
			cmd.Context(),
			parsedConfig,
			app,
			metrics,
		)

		if err != nil {
			return err
		}

		metricsAddress := fmt.Sprintf("%s:%d", cfg.Metrics.Host, cfg.Metrics.Port)

		m.Start(metricsAddress, metrics.Registry)

		// TODO: Add signal handling and gracefull shutdown
		return srv.Start()
	},
}
