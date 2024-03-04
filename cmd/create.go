package cmd

import (
	"fmt"

	go_bagit "github.com/nyudlts/go-bagit"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.Flags().StringVar(&inputDir, "input-dir", "", "the directory to be bagged")
	createCmd.Flags().StringVar(&checksumAlgorithm, "algorithm", "md5", "the algorithm used for checksums")
	createCmd.Flags().IntVar(&numProcesses, "processes", 1, "Use multiple processes to calculate checksums faster (default: 1)")
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		bag, err := go_bagit.CreateBag(inputDir, checksumAlgorithm, numProcesses)
		if err != nil {
			panic(err)
		}

		fmt.Println(bag.Name)
	},
}
