package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
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
        AccessKeySecret: utils.GetEnvVariable("AWS_SECRET_KY"),
        Bucket:          utils.GetEnvVariable("AWS_BUCKET"),
    }

    sess := utils.CreateSession(awsConfig)

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Unable to read file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    tempFilePath := "/tmp/" + header.Filename
    tempFile, err := os.Create(tempFilePath)
    if err != nil {
        http.Error(w, "Unable to create temp file", http.StatusInternalServerError)
        return
    }
    defer tempFile.Close()

    _, err = tempFile.ReadFrom(file)
    if err != nil {
        http.Error(w, "Unable to save temp file", http.StatusInternalServerError)
        return
    }

    err = utils.UploadObject(awsConfig.Bucket, tempFilePath, header.Filename, sess, awsConfig)
    if err != nil {
        http.Error(w, "Unable to upload file to S3", http.StatusInternalServerError)
        return
    }

    fileURL := utils.GetFileURL(awsConfig.Bucket, header.Filename, sess)
    fmt.Fprintf(w, "File uploaded successfully: %s", fileURL)
}