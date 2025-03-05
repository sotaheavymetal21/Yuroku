package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OnsenImage は温泉画像を表すエンティティです
type OnsenImage struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid"`
	OnsenID     string             `json:"onsen_id" bson:"onsen_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewOnsenImage は新しい温泉画像エンティティを作成します
func NewOnsenImage(onsenID, userID, imageURL, description string) *OnsenImage {
	now := time.Now()
	return &OnsenImage{
		UUID:        uuid.New().String(),
		OnsenID:     onsenID,
		UserID:      userID,
		ImageURL:    imageURL,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
