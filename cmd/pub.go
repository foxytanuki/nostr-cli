package cmd

import (
	"github.com/spf13/cobra"
)

var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Publish events",
	Long:  "Publish events",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
