package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
	"github.com/google/uuid"
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

func (h *DefaultHandler) UpdateUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    userID := vars["userID"]

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.SendErrorResponse(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	ext := filepath.Ext(header.Filename)
	uniqueFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	s3Key := fmt.Sprintf("profile-pictures/%s", uniqueFileName)

	// Start a transaction

	// Upload image to S3 bucket
	awsConfig := utils.GetAWSConfig()
	err = utils.UploadObject(awsConfig, s3Key, file)
	if err != nil {
		utils.SendErrorResponse(w, "Unable to upload file to S3", http.StatusInternalServerError)
		return
	}

	// Update image in db
	err = h.Storage.UpdateProfilePicture(userID, uniqueFileName)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to update user picture", http.StatusInternalServerError)
		return
	}


	presignedURL := utils.GetPresignURL(awsConfig, s3Key)

	response := struct {
		ImageURL string `json:"image_url"`
	}{
		ImageURL: presignedURL,
	}

	utils.SendDataResponse(w, "Profile picture updated successfully", http.StatusOK, response)
}