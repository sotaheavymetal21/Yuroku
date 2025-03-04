package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OnsenImage は温泉画像を表すエンティティです
type OnsenImage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID      string             `json:"uuid" bson:"uuid"`
	OnsenID   string             `json:"onsen_id" bson:"onsen_id"`
	ImageURL  string             `json:"image_url" bson:"image_url"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// NewOnsenImage は新しい温泉画像エンティティを作成します
func NewOnsenImage(onsenID, imageURL string) *OnsenImage {
	return &OnsenImage{
		UUID:      uuid.New().String(),
		OnsenID:   onsenID,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}
}
