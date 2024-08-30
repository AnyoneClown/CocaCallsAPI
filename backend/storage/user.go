package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AnyoneClown/CocaCallsAPI/types"
	"github.com/AnyoneClown/CocaCallsAPI/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (c *CockroachDB) CreateUser(user types.UserToCreate) (types.User, error) {
	if user.Provider == "" {
		if err := types.ValidateUser(user.Email, user.Password); err != nil {
			return types.User{}, err
		}
	} else {
		if user.GoogleID == "" || user.Email == "" {
			return types.User{}, fmt.Errorf("invalid OAuth user data")
		}
	}

	var existingUser types.User
	if err := c.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return types.User{}, fmt.Errorf("email already in use")
	} else if err != gorm.ErrRecordNotFound {
		return types.User{}, err
	}

	if user.Provider != "" {
		if err := c.DB.Where("google_id = ?", user.GoogleID).First(&existingUser).Error; err == nil {
			return types.User{}, fmt.Errorf("google ID already in use")
		} else if err != gorm.ErrRecordNotFound {
			return types.User{}, err
		}
	}

	userToCreate := types.User{
		ID:            uuid.New(),
		Email:         user.Email,
		Password:      user.Password,
		GoogleID:      user.GoogleID,
		Picture:       user.Picture,
		Provider:      user.Provider,
		IsAdmin:       user.IsAdmin,
		VerifiedEmail: user.VerifiedEmail,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	result := c.DB.Create(&userToCreate)
	if result.Error != nil {
		return types.User{}, result.Error
	}

	return userToCreate, nil
}

func (c *CockroachDB) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	result := c.DB.Preload("Subscription").Where("email = ?", email).First(&user)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		if result.Error == gorm.ErrRecordNotFound {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, result.Error
	}

	return user, nil
}

func (c *CockroachDB) GetUserByID(userID string) (types.User, error) {
	var user types.User
	result := c.DB.Preload("Subscription").Omit("password").Where("id = ?", userID).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, result.Error
	}

	awsConfig := utils.GetAWSConfig()
	folder := "profile-pictures/"
	if !strings.HasPrefix(user.Picture, folder) {
		user.Picture = folder + user.Picture
	}
	user.Picture = utils.GetPresignURL(awsConfig, user.Picture)
	log.Println("Presigned URL:", user.Picture)
	return user, nil
}

func (c *CockroachDB) UpdateUser(user *types.User) error {
	result := c.DB.Model(&types.User{}).
		Where("email = ?", user.Email).
		Updates(map[string]interface{}{
			"google_id":      user.GoogleID,
			"picture":        user.Picture,
			"provider":       user.Provider,
			"verified_email": user.VerifiedEmail,
			"updated_at":     time.Now(),
		})
	if result.Error != nil {
		log.Println("Failed to execute update statement:", result.Error)
		return result.Error
	}

	return nil
}

func (c *CockroachDB) GetUsers() ([]types.User, error) {
	var userResponses []types.User
	if err := c.DB.Preload("Subscription").
		Omit("password").
		Find(&userResponses).Error; err != nil {
		return nil, err
	}

	return userResponses, nil
}

func (c *CockroachDB) UpdateProfilePicture(id string, s3Key string) error {
	result := c.DB.Model(&types.User{}).
		Where("id = ?", id).
		Update("picture", s3Key)
	if result.Error != nil {
		log.Println("Failed to update profile picture:", result.Error)
		return result.Error
	}

	return nil
}
