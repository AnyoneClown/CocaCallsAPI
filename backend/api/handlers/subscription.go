package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/utils"
	"github.com/google/uuid"
)

type SubscriptionRequest struct {
	UserID uuid.UUID `json:"userID"`
}

func (h *DefaultHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var request SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
	}

	subscription, err := h.Storage.CreateSubscription(request.UserID)
	if err != nil {
		utils.SendErrorResponse(w, "Subsciption already exists", http.StatusBadRequest)
		return
	}

	utils.SendDataResponse(w, "Successfully created subscription", http.StatusCreated, subscription)
}
