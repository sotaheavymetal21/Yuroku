package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/yuroku/internal/domain/entity"
	"github.com/yourusername/yuroku/internal/infrastructure/database"
)

func main() {
	log.Println("シードデータの投入を開始します...")

	// MongoDBに接続
	db, err := database.GetMongoDB()
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}

	// コレクション設定
	usersCollection := db.Collection("users")
	onsenLogsCollection := db.Collection("onsen_logs")
	onsenImagesCollection := db.Collection("onsen_images")

	// 既存データのクリア（オプション）
	ctx := context.Background()
	log.Println("既存のコレクションをクリアします...")
	usersCollection.Drop(ctx)
	onsenLogsCollection.Drop(ctx)
	onsenImagesCollection.Drop(ctx)

	// ユーザーデータの作成
	users := createUsers()
	log.Printf("%d人のユーザーデータを作成しました", len(users))

	// ユーザーデータの挿入
	for _, user := range users {
		_, err := usersCollection.InsertOne(ctx, user)
		if err != nil {
			log.Printf("ユーザー挿入エラー: %v", err)
		}
	}
	log.Println("ユーザーデータの挿入が完了しました")

	// 温泉ログデータの作成
	onsenLogs := createOnsenLogs(users)
	log.Printf("%d件の温泉ログデータを作成しました", len(onsenLogs))

	// 温泉ログデータの挿入
	var insertedOnsenLogs []*entity.OnsenLog
	for _, onsenLog := range onsenLogs {
		result, err := onsenLogsCollection.InsertOne(ctx, onsenLog)
		if err != nil {
			log.Printf("温泉ログ挿入エラー: %v", err)
			continue
		}
		onsenLog.ID = result.InsertedID.(primitive.ObjectID)
		insertedOnsenLogs = append(insertedOnsenLogs, onsenLog)
	}
	log.Println("温泉ログデータの挿入が完了しました")

	// 温泉画像データの作成
	onsenImages := createOnsenImages(insertedOnsenLogs, users)
	log.Printf("%d件の温泉画像データを作成しました", len(onsenImages))

	// 温泉画像データの挿入
	for _, onsenImage := range onsenImages {
		_, err := onsenImagesCollection.InsertOne(ctx, onsenImage)
		if err != nil {
			log.Printf("温泉画像挿入エラー: %v", err)
		}
	}
	log.Println("温泉画像データの挿入が完了しました")

	log.Println("シードデータの投入が完了しました！")
}

// ユーザーデータを作成する関数
func createUsers() []*entity.User {
	// パスワードをハッシュ化する関数
	hashPassword := func(password string) string {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashed)
	}

	now := time.Now()
	defaultPassword := hashPassword("password123")

	return []*entity.User{
		{
			ID:        primitive.NewObjectID(),
			UUID:      uuid.New().String(),
			Name:      "田中太郎",
			Email:     "tanaka@example.com",
			Password:  defaultPassword,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        primitive.NewObjectID(),
			UUID:      uuid.New().String(),
			Name:      "鈴木花子",
			Email:     "suzuki@example.com",
			Password:  defaultPassword,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        primitive.NewObjectID(),
			UUID:      uuid.New().String(),
			Name:      "佐藤一郎",
			Email:     "sato@example.com",
			Password:  defaultPassword,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// 温泉ログデータを作成する関数
func createOnsenLogs(users []*entity.User) []*entity.OnsenLog {
	var onsenLogs []*entity.OnsenLog
	now := time.Now()

	// ユーザーごとに温泉ログを作成
	for _, user := range users {
		// 田中太郎の温泉ログ
		if user.Name == "田中太郎" {
			onsenLogs = append(onsenLogs, &entity.OnsenLog{
				ID:         primitive.NewObjectID(),
				UUID:       uuid.New().String(),
				UserID:     user.UUID,
				Name:       "草津温泉",
				Location:   "群馬県吾妻郡草津町",
				SpringType: entity.SpringTypeSulfur,
				Features: []entity.Feature{
					entity.FeatureOutdoorBath,
					entity.FeatureViewpoint,
					entity.FeatureHistorical,
				},
				VisitDate: time.Now().AddDate(0, -1, 0),
				Rating:    5,
				Comment:   "湯畑が素晴らしく、お湯の質も最高でした。硫黄の香りが強いですが、肌に良い感じがします。次回はゆっくり宿泊したいです。",
				CreatedAt: now,
				UpdatedAt: now,
			})

			onsenLogs = append(onsenLogs, &entity.OnsenLog{
				ID:         primitive.NewObjectID(),
				UUID:       uuid.New().String(),
				UserID:     user.UUID,
				Name:       "箱根湯本温泉",
				Location:   "神奈川県足柄下郡箱根町",
				SpringType: entity.SpringTypeSalt,
				Features: []entity.Feature{
					entity.FeatureOutdoorBath,
					entity.FeatureRestaurant,
					entity.FeatureAccommodation,
				},
				VisitDate: time.Now().AddDate(0, -3, 0),
				Rating:    4,
				Comment:   "都心から近く便利。塩化物泉で湯冷めしにくいです。箱根の自然も楽しめて、温泉街の雰囲気も良かったです。",
				CreatedAt: now,
				UpdatedAt: now,
			})
		}

		// 鈴木花子の温泉ログ
		if user.Name == "鈴木花子" {
			onsenLogs = append(onsenLogs, &entity.OnsenLog{
				ID:         primitive.NewObjectID(),
				UUID:       uuid.New().String(),
				UserID:     user.UUID,
				Name:       "城崎温泉",
				Location:   "兵庫県豊岡市城崎町",
				SpringType: entity.SpringTypeSimple,
				Features: []entity.Feature{
					entity.FeatureOutdoorBath,
					entity.FeatureHistorical,
					entity.FeatureAccommodation,
				},
				VisitDate: time.Now().AddDate(0, -2, 0),
				Rating:    5,
				Comment:   "風情ある温泉街が素敵。浴衣で外湯めぐりを楽しめました。7つの外湯を全て巡り、それぞれに特徴があって面白かったです。",
				CreatedAt: now,
				UpdatedAt: now,
			})
		}

		// 佐藤一郎の温泉ログ
		if user.Name == "佐藤一郎" {
			onsenLogs = append(onsenLogs, &entity.OnsenLog{
				ID:         primitive.NewObjectID(),
				UUID:       uuid.New().String(),
				UserID:     user.UUID,
				Name:       "別府温泉",
				Location:   "大分県別府市",
				SpringType: entity.SpringTypeSulfur,
				Features: []entity.Feature{
					entity.FeatureOutdoorBath,
					entity.FeatureDirectFromSpring,
				},
				VisitDate: time.Now().AddDate(0, -1, -15),
				Rating:    4,
				Comment:   "地獄めぐりが面白かった。多様な泉質が楽しめる温泉郷です。特に砂湯が気持ち良かったです。",
				CreatedAt: now,
				UpdatedAt: now,
			})

			onsenLogs = append(onsenLogs, &entity.OnsenLog{
				ID:         primitive.NewObjectID(),
				UUID:       uuid.New().String(),
				UserID:     user.UUID,
				Name:       "黒川温泉",
				Location:   "熊本県南小国町",
				SpringType: entity.SpringTypeSimple,
				Features: []entity.Feature{
					entity.FeatureOutdoorBath,
					entity.FeaturePrivateBath,
					entity.FeatureViewpoint,
				},
				VisitDate: time.Now().AddDate(0, -4, 0),
				Rating:    5,
				Comment:   "自然に囲まれた露天風呂が最高。入湯手形で3つの温泉を巡りました。周囲の自然と調和した温泉で癒されました。",
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
	}

	return onsenLogs
}

// 温泉画像データを作成する関数
func createOnsenImages(onsenLogs []*entity.OnsenLog, users []*entity.User) []*entity.OnsenImage {
	var onsenImages []*entity.OnsenImage
	now := time.Now()

	// 各温泉ログに画像を追加
	for _, onsenLog := range onsenLogs {
		// 基本的な画像データ
		baseImage := &entity.OnsenImage{
			ID:          primitive.NewObjectID(),
			UUID:        uuid.New().String(),
			OnsenID:     onsenLog.UUID,
			UserID:      onsenLog.UserID,
			ImageURL:    "/uploads/sample/" + onsenLog.Name + "_1.jpg", // 実際のパスに合わせて調整
			Description: onsenLog.Name + "の外観",
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		onsenImages = append(onsenImages, baseImage)

		// 一部の温泉には2枚目の画像も追加
		if onsenLog.Name == "草津温泉" || onsenLog.Name == "城崎温泉" || onsenLog.Name == "黒川温泉" {
			secondImage := &entity.OnsenImage{
				ID:          primitive.NewObjectID(),
				UUID:        uuid.New().String(),
				OnsenID:     onsenLog.UUID,
				UserID:      onsenLog.UserID,
				ImageURL:    "/uploads/sample/" + onsenLog.Name + "_2.jpg", // 実際のパスに合わせて調整
				Description: onsenLog.Name + "の露天風呂",
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			onsenImages = append(onsenImages, secondImage)
		}
	}

	return onsenImages
}
