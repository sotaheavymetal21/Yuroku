package gateway

import (
	"context"
	"io"

	"github.com/yourusername/yuroku/internal/domain/repository"
	"github.com/yourusername/yuroku/internal/infrastructure/storage"
)

// LocalStorageRepository はLocalFileStorageをStorageRepositoryインターフェースに適合させるアダプターです
type LocalStorageRepository struct {
	storage *storage.LocalFileStorage
}

// NewLocalStorageRepository は新しいLocalStorageRepositoryを作成します
func NewLocalStorageRepository(storage *storage.LocalFileStorage) *LocalStorageRepository {
	return &LocalStorageRepository{
		storage: storage,
	}
}

// Upload はファイルをアップロードします
func (r *LocalStorageRepository) Upload(ctx context.Context, file io.Reader, fileName, contentType string) (string, error) {
	return r.storage.Upload(ctx, file, fileName, contentType)
}

// Delete はファイルを削除します
func (r *LocalStorageRepository) Delete(ctx context.Context, fileURL string) error {
	return r.storage.Delete(ctx, fileURL)
}

// Ensure LocalStorageRepository implements StorageRepository
var _ repository.StorageRepository = (*LocalStorageRepository)(nil)
