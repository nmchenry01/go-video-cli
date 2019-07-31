package cmd

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/nmchenry/go-video-cli/cmd/decrypt"
	"github.com/nmchenry/go-video-cli/cmd/encrypt"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var author string

func init() {
	rootCmd.AddCommand(encrypt.Encrypt)
	rootCmd.AddCommand(decrypt.Decrypt)
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

func createHash(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
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

func convertBytesToMb(bytes int) int {
	return bytes / 1024 / 1024
}
