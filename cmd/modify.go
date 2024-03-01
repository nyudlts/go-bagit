package cmd

import (
	"log"

	go_bagit "github.com/nyudlts/go-bagit"
	"github.com/spf13/cobra"
)

func init() {
	tagManifestCmd.Flags().StringVar(&bagLocation, "bag", "", "bag to be validated")
	tagManifestCmd.Flags().BoolVar(&addBagFile, "add-to-bag", false, "add a file to tag manifest")
	tagManifestCmd.Flags().StringVar(&file, "file", "", "location of a file")
	rootCmd.AddCommand(tagManifestCmd)
}

var tagManifestCmd = &cobra.Command{
	Use: "modify",
	Run: func(cmd *cobra.Command, args []string) {
		if addBagFile == true {
			//Add a file to a bag
			b, err := go_bagit.GetExistingBag(bagLocation)
			if err != nil {
				panic(err)
			}
			if err := b.AddFileToBagRoot(file); err != nil {
				panic(err)
			}
		} else {
			log.Println("- WARNING - No valid subcommand provided")
		}
	},
}
