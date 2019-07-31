package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	Decrypt.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli decrypt [FILEPATH] [KEYNAME] [REGION]\n")
}

var Decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file provided on the command line",
	Long:  "Decrypt a file provided on the command line by providing a path the file, the name of the key stored in AWS Secrets Manager, and the AWS region where the key is stored",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		// Arguments
		filepath := args[0]
		keyName := args[1]
		region := args[2]

		// Create AWS Session
		log.Info("Creating AWS Session")
		session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region)}),
		)

		// Create Secrets Manager service
		svc := secretsmanager.New(session)

		// Build Get Secrets Value input
		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(keyName),
		}

		// Get Secrets Value
		log.Info("Retrieving Secret Key")
		result, err := svc.GetSecretValue(input)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to get secret key")
		}

		// Read in file
		log.Info("Reading input file")
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to read input")
		}

		// Decrypting the data
		log.Info("Decrypting File")
		plaintext := createPlaintext(data, *result.SecretString)

		basePath := strings.Split(filepath, ".encrypted")[0]

		// Write the data out to a file
		log.Info(fmt.Sprintf("Outputting file to your directory as %s", basePath))
		err = ioutil.WriteFile(basePath, plaintext, 0777)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to output file")
		}

		log.Info("Success!")
	},
}

func createPlaintext(data []byte, key string) []byte {
	byteKey := []byte(key)

	block, err := aes.NewCipher(byteKey)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to create cipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to create GCM")
	}

	nonceSize := gcm.NonceSize()

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to get plaintext")
	}

	return plaintext
}
