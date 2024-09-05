package storage

import (
	"fmt"
	"time"

	"github.com/AnyoneClown/CocaCallsAPI/types"
	"github.com/google/uuid"
)

func (c *CockroachDB) CreateSubscription(userID uuid.UUID) (types.Subscription, error) {
	var subscription types.Subscription
	if err := c.DB.Where("userID = ?", userID).First(&subscription); err == nil {
		return types.Subscription{}, fmt.Errorf("subscription already exists")
	}

	subscriptionToCreate := types.Subscription{
		ID: uuid.New(),
		UserID: userID,
		StartDate: time.Now(),
		EndDate: time.Now().AddDate(0, 1, 1),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := c.DB.Create(&subscriptionToCreate)
	if result.Error != nil {
		return types.Subscription{}, result.Error
	}

	return subscriptionToCreate, nil
}
