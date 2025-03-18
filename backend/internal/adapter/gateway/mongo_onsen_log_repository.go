package gateway

import (
	"context"
	"errors"
	"log"
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

// コレクション名とインデックス名
const (
	onsenLogsCollection = "onsen_logs"
	userIDIndex         = "user_id_idx"
	visitDateIndex      = "visit_date_idx"
	userVisitDateIndex  = "user_visit_date_idx"
	userSpringTypeIndex = "user_spring_type_idx"
	userLocationIndex   = "user_location_idx"
	userRatingIndex     = "user_rating_idx"
	compoundFilterIndex = "user_filter_compound_idx"
)

// NewMongoOnsenLogRepository は新しいMongoDBの温泉メモリポジトリを作成します
func NewMongoOnsenLogRepository(db *mongo.Database) *MongoOnsenLogRepository {
	// リポジトリインスタンスを作成
	repo := &MongoOnsenLogRepository{
		collection: db.Collection(onsenLogsCollection),
	}

	// 必要なインデックスを初期化
	go repo.ensureIndexes(context.Background())

	return repo
}

// ensureIndexes は必要なインデックスを設定します
func (r *MongoOnsenLogRepository) ensureIndexes(ctx context.Context) {
	// コンテキストをキャンセル可能にする
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// ユーザーIDのインデックス
	userIDIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetName(userIDIndex),
	}

	// 訪問日のインデックス
	visitDateIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "visit_date", Value: -1}},
		Options: options.Index().SetName(visitDateIndex),
	}

	// ユーザーID+訪問日の複合インデックス（よく使われる検索パターン）
	userVisitDateIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "visit_date", Value: -1}},
		Options: options.Index().SetName(userVisitDateIndex),
	}

	// ユーザーID+泉質の複合インデックス
	userSpringTypeIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "spring_type", Value: 1}},
		Options: options.Index().SetName(userSpringTypeIndex),
	}

	// ユーザーID+所在地の複合インデックス（テキスト検索の高速化）
	userLocationIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "location", Value: "text"}},
		Options: options.Index().SetName(userLocationIndex),
	}

	// ユーザーID+評価の複合インデックス
	userRatingIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "rating", Value: -1}},
		Options: options.Index().SetName(userRatingIndex),
	}

	// フィルタリングに使われる複合インデックス
	filterIdx := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "spring_type", Value: 1},
			{Key: "rating", Value: -1},
			{Key: "visit_date", Value: -1},
		},
		Options: options.Index().SetName(compoundFilterIndex),
	}

	// すべてのインデックスを一括で作成（存在する場合は無視される）
	indexes := []mongo.IndexModel{
		userIDIdx,
		visitDateIdx,
		userVisitDateIdx,
		userSpringTypeIdx,
		userLocationIdx,
		userRatingIdx,
		filterIdx,
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("Failed to create indexes: %v", err)
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
	// タイムアウト設定
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 検索条件を作成
	filter := bson.M{"user_id": userID}

	// パイプラインを使って効率的にクエリを実行
	pipeline := mongo.Pipeline{
		// マッチングステージ
		{{"$match", filter}},
		// ソートステージ
		{{"$sort", bson.M{"visit_date": -1}}},
		// ファセット（集計）ステージ - 一度のクエリで合計カウントとページングデータを取得
		{{"$facet", bson.M{
			"metadata": mongo.Pipeline{
				{{"$count", "total"}},
			},
			"data": mongo.Pipeline{
				{{"$skip", (page - 1) * limit}},
				{{"$limit", limit}},
			},
		}}},
	}

	// パイプラインを実行
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 結果構造体
	var results []struct {
		Metadata []struct {
			Total int `bson:"total"`
		} `bson:"metadata"`
		Data []*entity.OnsenLog `bson:"data"`
	}

	// 結果を取得
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	// 結果が空の場合
	if len(results) == 0 {
		return []*entity.OnsenLog{}, 0, nil
	}

	// カウントを取得
	totalCount := 0
	if len(results[0].Metadata) > 0 {
		totalCount = results[0].Metadata[0].Total
	}

	return results[0].Data, totalCount, nil
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
