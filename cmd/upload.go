package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cheggaaa/pb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upload)
	upload.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli upload [FILEPATH] [S3BUCKET] [REGION] [KEYNAME]\n")
}

var upload = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to an S3 bucket",
	Long:  "Upload a file to an S3 bucket given the path to the file, the bucket name, a key name",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		filepath := args[0]
		bucket := args[1]
		region := args[2]
		key := args[3]

		concurrency := 5     // Set max GoRoutines for upload to 5
		partSize := int64(5) // Upload in 5mb chunks

		// Create AWS Session
		log.Info("Creating AWS Session")
		session := createAWSSession(region)

		// Load File to upload
		log.Info("Reading File")
		file, bar := readFile(filepath)

		// Create upload manager for concurrent object upload
		log.Info("Creating Upload Manager")
		uploader := createUploadManager(session, partSize, concurrency)

		// Build input struct for upload
		log.Info("Building Input Parameters")
		input := buildUploadInput(bucket, key, file)

		//Upload to S3
		log.Info(fmt.Sprintf("Uploading Key: %s to S3 Bucket: %s", key, bucket))
		bar.Start()
		uploadToS3(uploader, input)
		bar.Finish()

		log.Info("Success")
	},
}

func createAWSSession(region string) *session.Session {
	validSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)}),
	)

	return validSession
}

func readFile(filepath string) (*customReader, *pb.ProgressBar) {
	file, err := os.Open(filepath)
	check(err, "There was an issue reading the file")

	fileInfo, err := file.Stat()
	check(err, "There was an issue retrieving file stats")

	limit := convertBytesToMb(int(fileInfo.Size()))
	bar := createProgressBar(limit)

	reader := &customReader{
		fp:       file,
		size:     fileInfo.Size(),
		progress: bar,
	}

	return reader, bar
}

func createUploadManager(session *session.Session, partSize int64, concurrency int) *s3manager.Uploader {
	uploader := s3manager.NewUploader(session, func(u *s3manager.Uploader) {
		u.PartSize = partSize * 1024 * 1024 // In chunks of 1 mb, min is 5 mb
		u.Concurrency = concurrency         // Number of GoRoutines to allow
	})

	return uploader
}

func buildUploadInput(bucket string, key string, file *customReader) *s3manager.UploadInput {
	input := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	}

	return input
}

func uploadToS3(uploader *s3manager.Uploader, input *s3manager.UploadInput) {
	_, err := uploader.Upload(input)
	check(err, "There was an issue uploading the file to S3")
}

func createProgressBar(fileSize int) *pb.ProgressBar {
	bar := pb.New(fileSize)
	return bar
}
