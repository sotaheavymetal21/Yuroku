package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SpringType は温泉の泉質を表す型です
type SpringType string

// 泉質の定数
const (
	SpringTypeSulfur   SpringType = "硫黄泉"
	SpringTypeCarbonic SpringType = "炭酸泉"
	SpringTypeAlkaline SpringType = "アルカリ泉"
	SpringTypeAcidic   SpringType = "酸性泉"
	SpringTypeSalt     SpringType = "塩化物泉"
	SpringTypeIron     SpringType = "鉄泉"
	SpringTypeRadium   SpringType = "ラジウム泉"
	SpringTypeSimple   SpringType = "単純温泉"
	SpringTypeOther    SpringType = "その他"
	SpringTypeUnknown  SpringType = "不明"
)

// Feature は温泉の特徴を表す型です
type Feature string

// 特徴の定数
const (
	FeatureOutdoorBath      Feature = "露天風呂あり"
	FeaturePrivateBath      Feature = "貸切風呂あり"
	FeatureDirectFromSpring Feature = "源泉掛け流し"
	FeatureSauna            Feature = "サウナあり"
	FeatureRestaurant       Feature = "食事処あり"
	FeatureAccommodation    Feature = "宿泊施設あり"
	FeatureViewpoint        Feature = "景色が良い"
	FeatureHistorical       Feature = "歴史ある温泉"
)

// OnsenLog は温泉メモを表すエンティティです
type OnsenLog struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID       string             `json:"uuid" bson:"uuid"`
	UserID     string             `json:"user_id" bson:"user_id"`
	Name       string             `json:"name" bson:"name"`
	Location   string             `json:"location" bson:"location"`
	SpringType SpringType         `json:"spring_type" bson:"spring_type"`
	Features   []Feature          `json:"features" bson:"features"`
	VisitDate  time.Time          `json:"visit_date" bson:"visit_date"`
	Rating     int                `json:"rating" bson:"rating"`
	Comment    string             `json:"comment" bson:"comment"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewOnsenLog は新しい温泉メモエンティティを作成します
func NewOnsenLog(userID, name, location string, springType SpringType, features []Feature, visitDate time.Time, rating int, comment string) *OnsenLog {
	now := time.Now()
	return &OnsenLog{
		UUID:       uuid.New().String(),
		UserID:     userID,
		Name:       name,
		Location:   location,
		SpringType: springType,
		Features:   features,
		VisitDate:  visitDate,
		Rating:     rating,
		Comment:    comment,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Update は温泉メモの情報を更新します
func (o *OnsenLog) Update(name, location string, springType SpringType, features []Feature, visitDate time.Time, rating int, comment string) {
	o.Name = name
	o.Location = location
	o.SpringType = springType
	o.Features = features
	o.VisitDate = visitDate
	o.Rating = rating
	o.Comment = comment
	o.UpdatedAt = time.Now()
}

// ValidateRating は評価値が有効かどうかを検証します
func ValidateRating(rating int) bool {
	return rating >= 0 && rating <= 5
}
