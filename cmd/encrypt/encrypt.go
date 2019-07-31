package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

/*
	TODO Better documentation around the usage
	TODO Look in to more flags/validation options
	TODO Write tests
	TODO Benchmarking
*/

func init() {
	Encrypt.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli encrypt [FILEPATH] [KEYNAME] [REGION] \n")
}

var Encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file provided on the command line",
	Long:  "Encrypt a file provided on the command line by providing a path to the file, the name of the key stored in AWS Secrets Manager, and the AWS region where the key is stored",
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
		input := secretsmanager.GetSecretValueInput{
			SecretId: aws.String(keyName),
		}

		// Get Secrets Value
		log.Info("Retrieving Secret Key")
		result, err := svc.GetSecretValue(&input)
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

		// Encrypting the data
		log.Info("Encrypting File")
		cipherText := createCipherText(data, *result.SecretString)

		outputPath := fmt.Sprintf("%s.encrypted", filepath)

		// Write the data out to a file
		log.Info(fmt.Sprintf("Outputting file to your directory as %s", outputPath))
		err = ioutil.WriteFile(outputPath, cipherText, 0777)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to output encrypted file")
		}

		log.Info("Success!")
	},
}

func createCipherText(data []byte, key string) []byte {
	byteKey := []byte(key)

	block, err := aes.NewCipher(byteKey)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to create block cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to create GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to read nonce")
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}
