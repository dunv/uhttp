/*
Copyright © 2023 Daniel Unverricht (daniel@unverricht.net)
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "uhttp_cli",
	Short: "A brief description of your application",
	Long:  `Command-line tools for uhttp`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
