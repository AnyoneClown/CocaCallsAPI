package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

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

func UploadObject(awsConfig AWSConfig, s3Key string, file multipart.File) error {
    cfg := CreateSession(awsConfig)
    uploader := manager.NewUploader(s3.NewFromConfig(cfg))

    _, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
        Bucket: aws.String(awsConfig.Bucket),
        Key:    aws.String(s3Key),
        Body:   file,
    })

    if err != nil {
        log.Printf("Error uploading to S3: %v", err)
        return fmt.Errorf("unable to upload file to S3: %v", err)
    }

    log.Printf("File uploaded successfully. Key: %s", s3Key)
    return nil
}

func GetPresignURL(awsConfig AWSConfig, s3Key string) string {
    cfg := CreateSession(awsConfig)
    s3client := s3.NewFromConfig(cfg)
    presignClient := s3.NewPresignClient(s3client)
    presignedUrl, err := presignClient.PresignGetObject(context.Background(),
        &s3.GetObjectInput{
            Bucket: aws.String(awsConfig.Bucket),
            Key:    aws.String(s3Key),
        },
        s3.WithPresignExpires(time.Minute*15))
    if err != nil {
        log.Printf("Error generating presigned URL: %v", err)
        return ""
    }
    return presignedUrl.URL
}