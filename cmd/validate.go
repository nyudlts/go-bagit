package cmd

import (
	go_bagit "github.com/nyudlts/go-bagit"
	"github.com/spf13/cobra"
)

var bagLocation string
var fast bool

func init() {
	validateCmd.Flags().StringVar(&bagLocation, "bag", "", "bag to be validated")
	validateCmd.MarkFlagRequired("bag")
	validateCmd.Flags().BoolVar(&fast, "fast", false, "perform fast validation, oxum only")
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {
		go_bagit.ValidateBag(bagLocation, fast)
	},
}
