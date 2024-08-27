package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
}

func GetAWSConfig() AWSConfig {
	return AWSConfig{
		Region:          GetEnvVariable("AWS_REGION"),
		AccessKeyID:     GetEnvVariable("AWS_ACCESS_KEY"),
		AccessKeySecret: GetEnvVariable("AWS_SECRET_KEY"),
		Bucket:          GetEnvVariable("AWS_BUCKET"),
	}
}

func CreateSession(awsConfig AWSConfig) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsConfig.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsConfig.AccessKeyID,
			awsConfig.AccessKeySecret,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return cfg
}

func UploadObject(awsConfig AWSConfig, filename string, file *os.File) error {
	cfg := CreateSession(awsConfig)
	uploader := manager.NewUploader(s3.NewFromConfig(cfg))

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(awsConfig.Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		log.Printf("Error uploading to S3: %v", err)
		return fmt.Errorf("unable to upload file to S3: %v", err)
	}

	log.Printf("File uploaded successfully. URL: %s", result.Location)
	return nil
}
