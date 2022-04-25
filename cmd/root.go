package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "buaa-login",
	Short: "Use the command line to login the buaa campus network.",
	Long: `Use the command line to login the buaa campus network.
It mainly needs to implement several encryption algorithms,
and obtain the current IP and Token through the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
