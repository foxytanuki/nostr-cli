package cmd

import (
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nostr-cli",
	Short: "A CLI tool for nostrich",
	Long:  "A CLI tool for nostrich",
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate private and public keys",
	Long:  "Generate private and public keys and display in the form of hex and bech32",
	Args:  cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		sk := nostr.GeneratePrivateKey()
		nsec, err := nip19.EncodePrivateKey(sk)
		if err != nil {
			return err
		}
		pk, err := nostr.GetPublicKey(sk)
		if err != nil {
			return err
		}
		npub, err := nip19.EncodePublicKey(pk)
		if err != nil {
			return err
		}
		outputs := []string{
			"===Private Key (Do NOT share to anyone)===",
			sk,
			nsec,
			"",
			"===Public Key (Share with your friends!)===",
			pk,
			npub,
		}

		for _, s := range outputs {
			fmt.Println(s)
		}

		return nil
	},
}

func Execute() {
	rootCmd.AddCommand(genCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
