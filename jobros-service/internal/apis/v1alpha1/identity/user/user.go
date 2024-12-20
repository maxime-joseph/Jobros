package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system.
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email" binding:"required,email"`
	Phone        string             `json:"phoneNumber" bson:"phoneNumber" binding:"required"`
	Role         string             `json:"roleRef" bson:"roleRef" binding:"required"`
	Status       string             `json:"status" bson:"status" binding:"required"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	Verification struct {
		Email        bool `json:"email" bson:"email"`
		Phone        bool `json:"phone" bson:"phone"`
		Identity     bool `json:"identity" bson:"identity"`
		Professional bool `json:"professional" bson:"professional"`
	} `json:"verificationStatus" bson:"verificationStatus"`
	Security struct {
		MFAEnabled      bool      `json:"mfaEnabled" bson:"mfaEnabled"`
		LoginAttempts   int       `json:"loginAttempts" bson:"loginAttempts"`
		LastLogin       time.Time `json:"lastLogin" bson:"lastLogin"`
		LastUpdated     time.Time `json:"lastUpdated" bson:"lastUpdated"`
		PasswordChanged time.Time `json:"lastPasswordChange" bson:"lastPasswordChange"`
	} `json:"security" bson:"security"`
}
