/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"k2a/internal/k2a"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var config = k2a.K2AConfig{}

var k2aCmd = &cobra.Command{
	Use:   "k2a",
	Short: "Export an AsyncAPI specification",
	Long:  `Export an AsyncAPI specification for a Kafka cluster and Schema Registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		zap.L().Info("cli config", zap.Any("config", config))
		yaml, err := k2a.ExportAsyncApi(&config)
		if err != nil {
			zap.L().Warn("run error", zap.String("export error", err.Error()))
			return
		}
		os.WriteFile(config.File, yaml, 0644)
	},
	Example: `
# no auth, local kafka, local registry
cli k2a --topics demo,sample
# no auth
cli k2a --kurl prod.kafka.com --rurl http://prod.schema-registry.com --topics demo,sample
# for SASL_PLAINTEXT
cli k2a --kurl prod.kafka.com --rurl http://prod.schema-registry.com --topics demo --username admin --username admin-secret
# SASL_SSL
...
	`,
}

func init() {
	rootCmd.AddCommand(k2aCmd)

	k2aCmd.Flags().StringVar(&config.Topics, "topics", "", "Topics to export")
	k2aCmd.Flags().StringVar(&config.File, "file", "k2a.yaml", "Output file name")

	k2aCmd.Flags().StringVar(&config.Certificate, "cert", "", "The optional certificate file for client authentication")
	k2aCmd.Flags().StringVar(&config.KeyFile, "key-file", "", "The optional key file for client authentication")
	k2aCmd.Flags().StringVar(&config.CAFile, "ca-file", "", "The optional certificate authority file for TLS client authentication")
	k2aCmd.Flags().BoolVar(&config.TLSSkipVerify, "tls-skip-verify", true, "Whether to skip TLS server cert verification")
	k2aCmd.Flags().BoolVar(&config.UseTLS, "use-tls", false, "Use TLS to communicate with the kafka cluster")
}
