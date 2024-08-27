package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
)

func (h *DefaultHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Storage.GetUsers()
	if err != nil {
		utils.SendErrorResponse(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	utils.SendDataResponse(w, "Users retrieved successfully", http.StatusOK, users)

}

func (h *DefaultHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := h.Storage.GetUserByID(vars["userID"])
	if err != nil {
		utils.SendErrorResponse(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	utils.SendDataResponse(w, "Users retrieved successfully", http.StatusOK, user)
}

func (h *DefaultHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
    awsConfig := utils.AWSConfig{
        Region:          utils.GetEnvVariable("AWS_REGION"),
        AccessKeyID:     utils.GetEnvVariable("AWS_ACCESS_KEY"),
        AccessKeySecret: utils.GetEnvVariable("AWS_SECRET_KEY"),
        Bucket:          utils.GetEnvVariable("AWS_BUCKET"),
    }

    log.Printf("AWS Config: Region: %s, Bucket: %s", awsConfig.Region, awsConfig.Bucket)

    sess:= utils.CreateSession(awsConfig)

    file, header, err := r.FormFile("file")
    if err != nil {
        log.Printf("Error reading file: %v", err)
        http.Error(w, "Unable to read file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    filename := filepath.Base(header.Filename)
    log.Printf("Attempting to upload file: %s", filename)

    uploader := s3manager.NewUploader(sess)
    result, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(awsConfig.Bucket),
        Key:    aws.String(filename),
        Body:   file,
    })

    if err != nil {
        log.Printf("Error uploading to S3: %v", err)
        if aerr, ok := err.(awserr.Error); ok {
            log.Printf("AWS Error Code: %s, Message: %s", aerr.Code(), aerr.Message())
        }
        http.Error(w, "Unable to upload file to S3", http.StatusInternalServerError)
        return
    }

    log.Printf("File uploaded successfully. URL: %s", result.Location)
    fmt.Fprintf(w, "File uploaded successfully: %s", result.Location)
}