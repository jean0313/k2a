/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"k2a/cmd"
	"k2a/internal/k2a"

	"go.uber.org/zap"
)

var Version = "0.0.1"

func main() {
	zap.L().Info("start", zap.String("version", Version))
	cmd.Execute()
}

func init() {
	k2a.InitLog()
}
