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

var port string

var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "Start web server",
	Long:  `Start web server to export topics`,
	Run: func(cmd *cobra.Command, args []string) {
		fs, _ := fs.Sub(pub, "public")
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(fs))))
		http.HandleFunc("/export", export)

		zap.L().Info(fmt.Sprintf("Starting server listen on %s, http://localhost:%s\n", port, port))
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
			log.Fatal(err)
		}
	},
	Example: `
# no auth, local kafka, local schema registry
cli ws
# for SASL_PLAINTEXT
cli ws --kurl prod.kafka.com --rurl http://prod.schema-registry.com --username admin --password admin-secret
# SASL_SSL
...
	`,
}

func init() {
	rootCmd.AddCommand(wsCmd)

	wsCmd.Flags().StringVar(&port, "port", "8080", "server port to listen")
}

func export(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v\n", err), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &config)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v\n", err), http.StatusBadRequest)
		return
	}

	zap.L().Info("config", zap.Any("config", maskPassword(config)))
	if config.Topics == "" {
		http.Error(w, "error: topics should not be empty", http.StatusBadRequest)
		return
	}

	err = config.InitTopicOperations()
	if err != nil {
		zap.L().Error("param error", zap.String("param error", err.Error()))
		return
	}
	yaml, err := k2a.ExportAsyncApi(&config)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v\n", err.Error()), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(yaml))
}
