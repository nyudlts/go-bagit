package cmd

import (
	go_bagit "github.com/nyudlts/go-bagit"
	"github.com/spf13/cobra"
)

func init() {
	validateCmd.Flags().StringVar(&bagLocation, "bag", "", "bag to be validated")
	validateCmd.MarkFlagRequired("bag")
	validateCmd.Flags().BoolVar(&complete, "completeness-only", false, "Only test whether the bag directory has the number of files and total size specified in Payload-Oxum without performing checksum validation to detect corruption.")
	validateCmd.Flags().BoolVar(&fast, "fast", false, "Test whether the bag directory has the expected payload specified in the checksum manifests without performing checksum validation to detect corruption.")
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {
		go_bagit.ValidateBag(bagLocation, fast, complete)
	},
}
