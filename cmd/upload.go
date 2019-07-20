package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upload)
	upload.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli upload [FILEPATH] [S3BUCKET] [KEYNAME]\n")
}

var upload = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to an S3 bucket",
	Long:  "Upload a file to an S3 bucket given the path to the file, the bucket name, a key name",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}
