package cmd

import (
	go_bagit "github.com/nyudlts/go-bagit"
	"github.com/spf13/cobra"
)

var bagLocation string

func init() {
	validateCmd.Flags().StringVar(&bagLocation, "bag", "", "bag to be validated")
	validateCmd.MarkFlagRequired("bag")
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {
		go_bagit.ValidateBag(bagLocation)
	},
}
