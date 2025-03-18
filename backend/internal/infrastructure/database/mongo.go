package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	defaultMaxPoolSize     = 100
	defaultConnectTimeout  = 30 * time.Second
	defaultMaxConnIdleTime = 60 * time.Second
	maxRetries             = 5
	retryInterval          = 3 * time.Second
)

// NewMongoClient はMongoDBクライアントを作成します
func NewMongoClient(uri string) (*mongo.Client, error) {
	// 接続オプションを設定
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(defaultMaxPoolSize).
		SetConnectTimeout(defaultConnectTimeout).
		SetMaxConnIdleTime(defaultMaxConnIdleTime).
		SetRetryWrites(true)

	// リトライループで接続を確立
	var client *mongo.Client
	var err error
	var connected bool

	for i := 0; i < maxRetries; i++ {
		// タイムアウト付きのコンテキストを設定
		ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
		defer cancel()

		// MongoDBに接続
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("MongoDB接続試行 %d/%d 失敗: %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		// 接続確認
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = client.Ping(pingCtx, readpref.Primary())
		pingCancel()

		if err != nil {
			log.Printf("MongoDB接続確認 %d/%d 失敗: %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		connected = true
		log.Printf("MongoDB接続に成功しました")
		break
	}

	if !connected {
		return nil, err
	}

	return client, nil
}

// GetMongoDB はMongoDBデータベースインスタンスを取得します
func GetMongoDB() (*mongo.Database, error) {
	// 環境変数からURI取得
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://mongo:27017"
	}

	// クライアント作成
	client, err := NewMongoClient(uri)
	if err != nil {
		return nil, err
	}

	// データベース名取得
	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = "yuroku"
	}

	// データベースを返す
	return client.Database(dbName), nil
}
