package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient はMongoDBクライアントを作成します
func NewMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// クライアントオプションを設定
	clientOptions := options.Client().ApplyURI(uri)

	// MongoDBに接続
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 接続確認
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetMongoDB はMongoDBデータベースインスタンスを取得します
func GetMongoDB() (*mongo.Database, error) {
	// 環境変数からURI取得
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// クライアント作成
	client, err := NewMongoClient(uri)
	if err != nil {
		return nil, err
	}

	// データベース名取得
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "yuroku"
	}

	// データベースを返す
	return client.Database(dbName), nil
}
