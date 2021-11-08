package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use: "go-bagit",
}

//flag variables
var (
	addTagFile        bool
	bagLocation       string
	checksumAlgorithm string
	complete          bool
	fast              bool
	file              string
	inputDir          string
	numProcesses      int
)
