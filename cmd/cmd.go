package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nostr-cli",
	Short: "A CLI tool for nostrich",
	Long:  "A CLI tool for nostrich",
}

func Execute() {
	rootCmd.AddCommand(genCmd)
	initPubCmd()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
