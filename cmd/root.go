/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"k2a/internal/k2a"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Debug bool

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "root command",
	Long:  `root command`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringVar(&config.KafkaUrl, "kurl", k2a.DEFAULT_KAFKA_URL, "Kafka cluster broker url")
	viper.BindPFlag("kurl", rootCmd.PersistentFlags().Lookup("kurl"))

	rootCmd.PersistentFlags().StringVar(&config.SchemaRegistryUrl, "rurl", k2a.DEFAULT_SCHEMA_REGISTRY_URL, "Schema registry url")
	viper.BindPFlag("rurl", rootCmd.PersistentFlags().Lookup("rurl"))

	rootCmd.PersistentFlags().StringVarP(&config.UserName, "username", "u", "", "username for kafka sasl_plaintext auth")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))

	rootCmd.PersistentFlags().StringVarP(&config.Password, "password", "p", "", "password for kafka sasl_plaintext auth")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	rootCmd.PersistentFlags().StringVar(&config.Password, "spec-version", "1.0.0", "Version number of the output file.")
	viper.BindPFlag("spec-version", rootCmd.PersistentFlags().Lookup("spec-version"))
}
