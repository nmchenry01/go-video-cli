package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decrypt)
}

var decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file provided on the command line",
	Long:  "Decrypt a file provided on the command line",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}
