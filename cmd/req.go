package cmd

import "github.com/spf13/cobra"

func req(relay string) error {
	return nil
}

func initReqCmd(rootCmd *cobra.Command) {
	reqCmd := &cobra.Command{
		Use:   "req",
		Short: "Request events",
		Long:  "Request events",
		RunE: func(cmd *cobra.Command, args []string) error {
			relay, err := cmd.Flags().GetString("relay")
			if err != nil {
				return err
			}
			if err := req(relay); err != nil {
				return err
			}
			return nil
		},
	}
	reqCmd.Flags().StringP("relay", "r", "", "relay URL to publish an event")
	reqCmd.MarkFlagRequired("relay")
	rootCmd.AddCommand(reqCmd)
}
