package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genKey)

	genKey.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli genKey [KEY] [KEYNAME] [REGION]\n")
}

var genKey = &cobra.Command{
	Use:   "genKey",
	Short: "Update the key being used to encrypt and decrypt files for the CLI",
	Long:  "Update the key being used to encrypt and decrypt files for the CLI by providing a key, the name of the key, and the region",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		key := args[0]
		keyName := args[1]
		region := args[2]

		// Create AWS Session
		log.Info("Creating Session")
		session := createAWSSession(region)

		// Hash the key
		log.Info("Hashing the Key")
		hashedKey := createHash(key)

		// Create Secrets Manager service
		svc := secretsmanager.New(session)

		// Generate Secret Input
		log.Info("Building Secrets Manager Input")
		input := buildSecretsInput(hashedKey, keyName)

		// Upload to Secrets Manager
		log.Info("Upload to Secrets Manager")
		uploadToSecretsManager(svc, input)

		log.Info("Success!")
	},
}

func buildSecretsInput(hashedKey string, keyName string) *secretsmanager.PutSecretValueInput {
	input := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(keyName),
		SecretString: aws.String(hashedKey),
	}

	return input
}

func uploadToSecretsManager(svc *secretsmanager.SecretsManager, input *secretsmanager.PutSecretValueInput) {
	_, err := svc.PutSecretValue(input)
	check(err, "There was an issue uploading to Secrets Manager")
}
