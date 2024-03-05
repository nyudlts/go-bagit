package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use: "info",
	Run: func(cmd *cobra.Command, args []string) {
		bi, _ := debug.ReadBuildInfo()
		fmt.Println("Build Info\n==========")
		fmt.Println(bi)
	},
}
