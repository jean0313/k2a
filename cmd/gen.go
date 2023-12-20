/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"k2a/internal/gen"
	"k2a/internal/k2a"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	gCtx gen.GlobalContext
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate code",
	Long:  `generate code for asyncapi spec`,
	Run: func(cmd *cobra.Command, args []string) {
		if Debug {
			k2a.LOG_LEVEL.SetLevel(zap.DebugLevel)
		}
		gen.Gen(&gCtx)
	},
	Example: `cli gen --asyncapi-file k2a.yaml --dest-dir output`,
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVar(&gCtx.Artifact, "artifact", "sample-app", "artifact id for maven project")
	genCmd.Flags().StringVar(&gCtx.Group, "group", "com.sample", "group id for maven project")
	genCmd.Flags().StringVar(&gCtx.PackageName, "package-name", "sample", "package name for maven project")
	genCmd.Flags().StringVar(&gCtx.Description, "description", "this is a sample app", "a description for maven project")
	genCmd.Flags().StringVar(&gCtx.ReleaseVersion, "release-version", "1.0.0", "parent project version for maven project")
	genCmd.Flags().StringVar(&gCtx.DestDir, "dest-dir", "output", "output dir for generation")
	genCmd.Flags().StringVar(&gCtx.AsyncAPIFile, "asyncapi-file", "", "asyncapi spec file for generation")
}
