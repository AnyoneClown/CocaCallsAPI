package handlers

import (
	"net/http"

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

func (h *DefaultHandler) UpdateUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	file, header, err := r.FormFile("file")
    if err != nil {
        utils.SendErrorResponse(w, "Failed to retrieve file", http.StatusBadRequest)
        return
    }
    defer file.Close()

	

	utils.SendDataResponse(w, "Users retrieved successfully", http.StatusOK, user)
}