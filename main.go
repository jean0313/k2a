/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"k2a/cmd"
	"k2a/internal/k2a"
)

func main() {
	cmd.Execute()
}

func init() {
	k2a.InitLog()
}
