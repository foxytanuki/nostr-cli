package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func initRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nostr-cli",
		Short: "A CLI tool for nostrich",
		Long:  "A CLI tool for nostrich",
	}
	return rootCmd
}

func Execute() {
	rootCmd := initRootCmd()
	initGenCmd(rootCmd)
	initPubCmd(rootCmd)
	initReqCmd(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
