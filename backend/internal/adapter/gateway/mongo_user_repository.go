package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/yuroku/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoUserRepository はMongoDBを使用したユーザーリポジトリの実装です
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository は新しいMongoDBユーザーリポジトリを作成します
func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

// Create は新しいユーザーを作成します
func (r *MongoUserRepository) Create(ctx context.Context, user *entity.User) error {
	// 既存のユーザーをチェック
	existingUser, _ := r.FindByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("このメールアドレスは既に登録されています")
	}

	// ドキュメントを作成
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// MongoDBに保存
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// 生成されたIDを設定
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

// FindByID はIDでユーザーを検索します
func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User

	// IDがObjectIDの場合
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		// ObjectIDで検索
		err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
		if err == nil {
			return &user, nil
		}
	}

	// UUIDで検索
	err = r.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ユーザーが見つかりません")
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail はメールアドレスでユーザーを検索します
func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ユーザーが見つかりません")
		}
		return nil, err
	}

	return &user, nil
}

// Update はユーザー情報を更新します
func (r *MongoUserRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()

	// MongoDBを更新
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete はユーザーを削除します
func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
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
		return errors.New("ユーザーが見つかりません")
	}

	return nil
}
