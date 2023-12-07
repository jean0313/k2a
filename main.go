/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"k2a/cmd"
	"k2a/internal/k2a"
)

var Version = "0.0.1"

func main() {
	fmt.Printf("version: %v\n", Version)
	cmd.Execute()
}

func init() {
	k2a.InitLog()
}
