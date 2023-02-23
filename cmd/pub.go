package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/cobra"
)

func pub(relayUrl string, sk string, content string) error {
	relay, err := nostr.RelayConnect(context.Background(), relayUrl)
	if err != nil {
		return cmdError(ErrRelayConnection, err.Error())
	}

	if len(sk) < 64 || len(sk) > 64 {
		return cmdError(ErrInvalidHashLen, fmt.Sprintf("invalid hash length: %d, should be 32-bytes lowercase hex-encoded private key", len(sk)))
	}
	pub, err := nostr.GetPublicKey(sk)
	if err != nil {
		return cmdError(ErrGetPublicKey, err.Error())
	}

	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      1,
		Tags:      nil,
		Content:   content,
	}
	if err := ev.Sign(sk); err != nil {
		return cmdError(ErrSignEvent, err.Error())
	}

	status := relay.Publish(context.Background(), ev)
	fmt.Println("published to ", relayUrl, status)
	if status == nostr.PublishStatusFailed {
		return errors.New("error: publish failed")
	}
	return nil
}

func initPubCmd(rootCmd *cobra.Command) {
	pubCmd := &cobra.Command{
		Use:   "pub",
		Short: "Publish events",
		Long:  "Publish events",
		RunE: func(cmd *cobra.Command, args []string) error {
			relay, err := cmd.Flags().GetString("relay")
			if err != nil {
				return err
			}
			sk, err := cmd.Flags().GetString("secret")
			if err != nil {
				return err
			}
			content, err := cmd.Flags().GetString("content")
			if err != nil {
				return err
			}
			if err := pub(relay, sk, content); err != nil {
				return err
			}
			return nil
		},
	}
	pubCmd.Flags().StringP("relay", "r", "", "relay URL to publish an event")
	pubCmd.Flags().StringP("secret", "s", "", "secret key to publish an event")
	pubCmd.Flags().StringP("content", "c", "", "content of an event")
	pubCmd.MarkFlagRequired("relay")
	pubCmd.MarkFlagRequired("secret")
	pubCmd.MarkFlagRequired("content")
	rootCmd.AddCommand(pubCmd)
}
