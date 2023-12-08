/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"k2a/internal/k2a"

	"github.com/spf13/cobra"
)

var config = k2a.K2AConfig{}

var k2aCmd = &cobra.Command{
	Use:   "k2a",
	Short: "Export an AsyncAPI specification",
	Long:  `Export an AsyncAPI specification for a Kafka cluster and Schema Registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := k2a.ExportAsyncApi(config)
		if err != nil {
			panic(err)
		}
	},
	Example: `cli k2a --kurl prod.kafka.com:9092 --rurl http://prod.schema-registry.com --topics demo,sample`,
}

func init() {
	rootCmd.AddCommand(k2aCmd)

	k2aCmd.Flags().StringVar(&config.KafkaUrl, "kurl", "localhost:9092", "Kafka cluster broker url")
	k2aCmd.Flags().StringVar(&config.SchemaRegistryUrl, "rurl", "http://localhost:8081", "Schema registry url")
	k2aCmd.Flags().StringVar(&config.Topics, "topics", "", "Topics to export")
	k2aCmd.Flags().StringVar(&config.File, "file", "k2a.yaml", "Output file name")
	k2aCmd.Flags().StringVar(&config.SpecVersion, "spec-version", "1.0.0", "Version number of the output file.")
	k2aCmd.Flags().StringVar(&config.Certificate, "cert", "", "The optional certificate file for client authentication")
	k2aCmd.Flags().StringVar(&config.KeyFile, "key-file", "", "The optional key file for client authentication")
	k2aCmd.Flags().StringVar(&config.CAFile, "ca-file", "", "The optional certificate authority file for TLS client authentication")
	k2aCmd.Flags().BoolVar(&config.TLSSkipVerify, "tls-skip-verify", true, "Whether to skip TLS server cert verification")
	k2aCmd.Flags().BoolVar(&config.UseTLS, "use-tls", false, "Use TLS to communicate with the kafka cluster")
}
