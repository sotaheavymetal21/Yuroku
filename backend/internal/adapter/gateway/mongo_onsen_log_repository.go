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

// MongoOnsenLogRepository はMongoDBを使用した温泉メモリポジトリの実装です
type MongoOnsenLogRepository struct {
	collection *mongo.Collection
}

// NewMongoOnsenLogRepository は新しいMongoDBの温泉メモリポジトリを作成します
func NewMongoOnsenLogRepository(db *mongo.Database) *MongoOnsenLogRepository {
	return &MongoOnsenLogRepository{
		collection: db.Collection("onsen_logs"),
	}
}

// Create は新しい温泉メモを作成します
func (r *MongoOnsenLogRepository) Create(ctx context.Context, onsenLog *entity.OnsenLog) error {
	// ドキュメントを作成
	now := time.Now()
	onsenLog.CreatedAt = now
	onsenLog.UpdatedAt = now

	// MongoDBに保存
	result, err := r.collection.InsertOne(ctx, onsenLog)
	if err != nil {
		return err
	}

	// 生成されたIDを設定
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		onsenLog.ID = oid
	}

	return nil
}

// FindByID はIDで温泉メモを検索します
func (r *MongoOnsenLogRepository) FindByID(ctx context.Context, id string) (*entity.OnsenLog, error) {
	var onsenLog entity.OnsenLog

	// IDがObjectIDの場合
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		// ObjectIDで検索
		err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&onsenLog)
		if err == nil {
			return &onsenLog, nil
		}
	}

	// UUIDで検索
	err = r.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&onsenLog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("温泉メモが見つかりません")
		}
		return nil, err
	}

	return &onsenLog, nil
}

// FindByUserID はユーザーIDに紐づく温泉メモを検索します
func (r *MongoOnsenLogRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.OnsenLog, error) {
	// 検索条件を作成
	filter := bson.M{"user_id": userID}

	// ソート条件を作成（訪問日の降順）
	opts := options.Find().SetSort(bson.M{"visit_date": -1})

	// 検索を実行
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 結果を取得
	var onsenLogs []*entity.OnsenLog
	if err := cursor.All(ctx, &onsenLogs); err != nil {
		return nil, err
	}

	return onsenLogs, nil
}

// FindByUserIDWithPagination はユーザーIDに紐づく温泉メモをページネーションで検索します
func (r *MongoOnsenLogRepository) FindByUserIDWithPagination(ctx context.Context, userID string, page, limit int) ([]*entity.OnsenLog, int, error) {
	// 検索条件を作成
	filter := bson.M{"user_id": userID}

	// 総件数を取得
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// ページネーション条件を作成
	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.M{"visit_date": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// 検索を実行
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 結果を取得
	var onsenLogs []*entity.OnsenLog
	if err := cursor.All(ctx, &onsenLogs); err != nil {
		return nil, 0, err
	}

	return onsenLogs, int(totalCount), nil
}

// FindByUserIDAndFilter はユーザーIDと条件に紐づく温泉メモを検索します
func (r *MongoOnsenLogRepository) FindByUserIDAndFilter(ctx context.Context, userID string, springType entity.SpringType, location string, minRating int, startDate, endDate *time.Time, page, limit int) ([]*entity.OnsenLog, int, error) {
	// 検索条件を作成
	filter := bson.M{"user_id": userID}

	// 泉質でフィルタリング
	if springType != "" {
		filter["spring_type"] = springType
	}

	// 所在地でフィルタリング
	if location != "" {
		filter["location"] = bson.M{"$regex": location, "$options": "i"}
	}

	// 評価でフィルタリング
	if minRating > 0 {
		filter["rating"] = bson.M{"$gte": minRating}
	}

	// 訪問日でフィルタリング
	if startDate != nil || endDate != nil {
		dateFilter := bson.M{}
		if startDate != nil {
			dateFilter["$gte"] = startDate
		}
		if endDate != nil {
			dateFilter["$lte"] = endDate
		}
		filter["visit_date"] = dateFilter
	}

	// 総件数を取得
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// ページネーション条件を作成
	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.M{"visit_date": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// 検索を実行
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 結果を取得
	var onsenLogs []*entity.OnsenLog
	if err := cursor.All(ctx, &onsenLogs); err != nil {
		return nil, 0, err
	}

	return onsenLogs, int(totalCount), nil
}

// Update は温泉メモを更新します
func (r *MongoOnsenLogRepository) Update(ctx context.Context, onsenLog *entity.OnsenLog) error {
	onsenLog.UpdatedAt = time.Now()

	// MongoDBを更新
	filter := bson.M{"_id": onsenLog.ID}
	update := bson.M{"$set": onsenLog}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete は温泉メモを削除します
func (r *MongoOnsenLogRepository) Delete(ctx context.Context, id string) error {
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
		return errors.New("温泉メモが見つかりません")
	}

	return nil
}

// DeleteByUserID はユーザーIDに紐づく温泉メモをすべて削除します
func (r *MongoOnsenLogRepository) DeleteByUserID(ctx context.Context, userID string) error {
	// ユーザーIDで削除
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
