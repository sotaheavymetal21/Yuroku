package repository

import (
	"context"
	"io"
)

// StorageRepository はファイルストレージを担当するインターフェースです
type StorageRepository interface {
	// Upload はファイルをアップロードします
	Upload(ctx context.Context, file io.Reader, fileName, contentType string) (string, error)

	// Delete はファイルを削除します
	Delete(ctx context.Context, fileURL string) error
}
