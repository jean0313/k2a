/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"k2a/internal/k2a"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

//go:embed public
var pub embed.FS

var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "Start web server",
	Long:  `Start web server to export topics`,
	Run: func(cmd *cobra.Command, args []string) {
		fs, _ := fs.Sub(pub, "public")
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(fs))))
		http.HandleFunc("/export", export)

		zap.L().Info("Starting server listen on 8080, http://localhost:8080\n")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(wsCmd)
}

func export(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v\n", err), http.StatusBadRequest)
		return
	}

	var c k2a.K2AConfig
	err = json.Unmarshal(body, &c)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v\n", err), http.StatusBadRequest)
		return
	}
	defaultValue(&c)
	zap.L().Info("config", zap.Any("config", c))

	if c.Topics == "" {
		http.Error(w, "error: topics should not be empty", http.StatusBadRequest)
		return
	}

	yaml, err := k2a.ExportAsyncApi(c)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v\n", err.Error()), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(yaml))
}

func defaultValue(config *k2a.K2AConfig) {
	if config.KafkaUrl == "" {
		config.KafkaUrl = KAFKA_URL
	}

	if config.SchemaRegistryUrl == "" {
		config.SchemaRegistryUrl = SCHEMA_REGISTRY_URL
	}
}
