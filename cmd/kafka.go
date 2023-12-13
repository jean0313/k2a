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
	Short: "kafka command",
	Long:  `kafka command`,
	Run: func(cmd *cobra.Command, args []string) {
		topics, err := queryTopics(query)
		if err != nil {
			zap.L().Warn("run error", zap.String("export error", err.Error()))
			return
		}
		for _, v := range topics {
			fmt.Println(v)
		}
	},
}

func queryTopics(query string) ([]string, error) {
	detail := k2a.CreateAccountDetails(&config)
	return detail.SearchTopics(query)
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	kafkaCmd.Flags().StringVar(&query, "query", "", "keyword in topic name")
}
