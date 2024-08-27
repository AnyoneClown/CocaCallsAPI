package utils

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSConfig struct {
    Region          string
    AccessKeyID     string
    AccessKeySecret string
    Bucket          string
}
func CreateSession(awsConfig AWSConfig) *session.Session {
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(awsConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.AccessKeyID,
				awsConfig.AccessKeySecret,
				"",
			),
		},
	))
	return sess
}

func CreateS3Session(sess *session.Session) *s3.S3 {
	return s3.New(sess)
}

func UploadObject(bucket string, filePath string, fileName string, sess *session.Session, awsConfig AWSConfig) error {

	// Open file to upload
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("unable to open file %q, %v", err)
		return err
	}
	defer file.Close()

	// Upload to s3
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		fmt.Errorf("failed to upload object, %v\n", err)
		return err
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucket)
	return nil
}

func GetFileURL(bucket, fileName string, sess *session.Session) string {
    svc := s3.New(sess)
    req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(fileName),
    })
    urlStr, err := req.Presign(15 * 60)
    if err != nil {
        fmt.Println("Failed to sign request", err)
        return ""
    }
    return urlStr
}
