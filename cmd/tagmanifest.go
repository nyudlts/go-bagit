package cmd

import (
	go_bagit "github.com/nyudlts/go-bagit"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	tagManifestCmd.Flags().StringVar(&bagLocation, "bag", "", "bag to be validated")
	tagManifestCmd.Flags().BoolVar(&addTagFile, "add-tag", false, "add a file to tag manifest")
	tagManifestCmd.Flags().StringVar(&file, "file", "", "location of file")
	rootCmd.AddCommand(tagManifestCmd)
}

var tagManifestCmd = &cobra.Command{
	Use: "tag",
	Run: func(cmd *cobra.Command, args []string) {
		if addTagFile == true {
			if err := go_bagit.AddFileToTagManifest(bagLocation, file); err != nil {
				panic(err)
			}
		} else {
			log.Println("-WARNING- No valid subcommand provided")
		}
	},
}
