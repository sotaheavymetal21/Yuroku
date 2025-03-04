package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/yuroku/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoOnsenImageRepository はMongoDBを使用した温泉画像リポジトリの実装です
type MongoOnsenImageRepository struct {
	collection *mongo.Collection
}

// NewMongoOnsenImageRepository は新しいMongoDBの温泉画像リポジトリを作成します
func NewMongoOnsenImageRepository(db *mongo.Database) *MongoOnsenImageRepository {
	return &MongoOnsenImageRepository{
		collection: db.Collection("onsen_images"),
	}
}

// Create は新しい温泉画像を作成します
func (r *MongoOnsenImageRepository) Create(ctx context.Context, image *entity.OnsenImage) error {
	// ドキュメントを作成
	now := time.Now()
	image.CreatedAt = now
	image.UpdatedAt = now

	// MongoDBに保存
	result, err := r.collection.InsertOne(ctx, image)
	if err != nil {
		return err
	}

	// 生成されたIDを設定
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		image.ID = oid
	}

	return nil
}

// FindByID はIDで温泉画像を検索します
func (r *MongoOnsenImageRepository) FindByID(ctx context.Context, id string) (*entity.OnsenImage, error) {
	var image entity.OnsenImage

	// IDがObjectIDの場合
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		// ObjectIDで検索
		err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&image)
		if err == nil {
			return &image, nil
		}
	}

	// UUIDで検索
	err = r.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("温泉画像が見つかりません")
		}
		return nil, err
	}

	return &image, nil
}

// FindByOnsenID は温泉IDに紐づく画像を検索します
func (r *MongoOnsenImageRepository) FindByOnsenID(ctx context.Context, onsenID string) ([]*entity.OnsenImage, error) {
	// 検索条件を作成
	filter := bson.M{"onsen_id": onsenID}

	// ソート条件を作成（作成日時の降順）
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	// 検索を実行
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 結果を取得
	var images []*entity.OnsenImage
	if err := cursor.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// FindByOnsenIDAndUserID は温泉IDとユーザーIDに紐づく画像を検索します
func (r *MongoOnsenImageRepository) FindByOnsenIDAndUserID(ctx context.Context, onsenID, userID string) ([]*entity.OnsenImage, error) {
	// 検索条件を作成
	filter := bson.M{
		"onsen_id": onsenID,
		"user_id":  userID,
	}

	// ソート条件を作成（作成日時の降順）
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	// 検索を実行
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 結果を取得
	var images []*entity.OnsenImage
	if err := cursor.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// Update は温泉画像を更新します
func (r *MongoOnsenImageRepository) Update(ctx context.Context, image *entity.OnsenImage) error {
	image.UpdatedAt = time.Now()

	// MongoDBを更新
	filter := bson.M{"_id": image.ID}
	update := bson.M{"$set": image}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete は温泉画像を削除します
func (r *MongoOnsenImageRepository) Delete(ctx context.Context, id string) error {
	// IDがObjectIDの場合
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		// ObjectIDで削除
		_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err == nil {
			return nil
		}
	}

	// UUIDで削除
	result, err := r.collection.DeleteOne(ctx, bson.M{"uuid": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("温泉画像が見つかりません")
	}

	return nil
}

// DeleteByOnsenID は温泉IDに紐づく画像をすべて削除します
func (r *MongoOnsenImageRepository) DeleteByOnsenID(ctx context.Context, onsenID string) error {
	// 温泉IDで削除
	_, err := r.collection.DeleteMany(ctx, bson.M{"onsen_id": onsenID})
	return err
}

// DeleteByUserID はユーザーIDに紐づく画像をすべて削除します
func (r *MongoOnsenImageRepository) DeleteByUserID(ctx context.Context, userID string) error {
	// ユーザーIDで削除
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
