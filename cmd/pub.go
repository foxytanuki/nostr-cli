package cmd

import (
	"github.com/spf13/cobra"
)

var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Publish events",
	Long:  "Publish events",
	RunE:  pub,
}

func pub(cmd *cobra.Command, args []string) error {
	cmd.Flags().GetString("relay")
	// relay, e := nostr.RelayConnect(context.Background(), url)
	return nil
}

func initPubCmd() {
	pubCmd.PersistentFlags().String("relay", "", "relay URL to publish an event")
	pubCmd.MarkPersistentFlagRequired("relay")
	rootCmd.AddCommand(pubCmd)
}
