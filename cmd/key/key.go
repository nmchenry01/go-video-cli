package key

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	GenKey.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli genKey [KEY] [KEYNAME] [REGION]\n")
}

var GenKey = &cobra.Command{
	Use:   "genKey",
	Short: "Update the key being used to encrypt and decrypt files for the CLI",
	Long:  "Update the key being used to encrypt and decrypt files for the CLI by providing a key, the name of the key, and the region",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		key := args[0]
		keyName := args[1]
		region := args[2]

		// Create AWS Session
		log.Info("Creating AWS Session")
		session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region)}),
		)

		// Hash the key
		log.Info("Hashing the Key")
		hashedKey := createHash(key)

		// Create Secrets Manager service
		svc := secretsmanager.New(session)

		// Generate Secret Input
		input := &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(keyName),
			SecretString: aws.String(hashedKey),
		}

		// Upload to Secrets Manager
		log.Info("Uploading Key to Secrets Manager")
		putSecret(svc, input)

		log.Info("Success!")
	},
}

func putSecret(svc *secretsmanager.SecretsManager, input *secretsmanager.PutSecretValueInput) {
	_, err := svc.PutSecretValue(input)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to upload secret")
	}
}

func createHash(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
