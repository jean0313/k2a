/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"k2a/internal/k2a"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var query string

var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "command to query topics by topic name",
	Long:  `command to query topics by topic name`,
	Run: func(cmd *cobra.Command, args []string) {
		if query == "" {
			zap.L().Warn("run error", zap.String("error", "query should not be empty!"))
			return
		}
		topics, err := queryTopics(query)
		if err != nil {
			zap.L().Warn("run error", zap.String("export error", err.Error()))
			return
		}
		for _, v := range topics {
			fmt.Println(v)
		}
	},
	Example: `
# no auth, local kafka, local registry
cli kafka --query test
# no auth
cli kafka --kurl prod.kafka.com --rurl http://prod.schema-registry.com --query test
# for SASL_PLAINTEXT
cli kafka --kurl prod.kafka.com --rurl http://prod.schema-registry.com --query test --username admin --username admin-secret
# SASL_SSL
...`,
}

func queryTopics(query string) ([]string, error) {
	detail := k2a.CreateAccountDetails(&config)
	return detail.SearchTopics(query)
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	kafkaCmd.Flags().StringVar(&query, "query", "", "keyword in topic name")
}
