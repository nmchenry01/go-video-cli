package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/nmchenry/go-video-cli/cmd/decrypt"
	"github.com/nmchenry/go-video-cli/cmd/key"

	"github.com/nmchenry/go-video-cli/cmd/encrypt"
	"github.com/nmchenry/go-video-cli/cmd/upload"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var author string

func init() {
	rootCmd.AddCommand(encrypt.Encrypt)
	rootCmd.AddCommand(decrypt.Decrypt)
	rootCmd.AddCommand(key.GenKey)
	rootCmd.AddCommand(upload.Upload)
}

var rootCmd = &cobra.Command{
	Use:   "go-video-cli",
	Short: "A simple cli for encrypting files, decrypting files, uploading files to S3, and deleting files from S3",
	Long:  "A simple cli for encrypting files, decrypting files, uploading files to S3, and deleting files from S3",
}

// Execute : The main function to execute the rootCmd
func Execute() {
	err := rootCmd.Execute()
	check(err, "There was an issue running the root command")
}

func check(err error, message string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(message)
	}
}

func buildGetSecretInput(keyName string) *secretsmanager.GetSecretValueInput {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(keyName),
	}
	return input
}

func getSecret(svc *secretsmanager.SecretsManager, input *secretsmanager.GetSecretValueInput) string {
	result, err := svc.GetSecretValue(input)
	check(err, "There was an issue retrieving the secret from Secrets Manager")

	return *result.SecretString
}
