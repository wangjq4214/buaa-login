package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the program version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v" + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
