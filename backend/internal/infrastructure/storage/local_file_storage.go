package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalFileStorage はローカルファイルシステムを使用したストレージ実装です
type LocalFileStorage struct {
	basePath string
}

// NewLocalFileStorage は新しいLocalFileStorageインスタンスを作成します
func NewLocalFileStorage(basePath string) (*LocalFileStorage, error) {
	// ベースディレクトリが存在しない場合は作成
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create storage directory: %w", err)
		}
	}

	return &LocalFileStorage{
		basePath: basePath,
	}, nil
}

// SaveFile はファイルを保存します
func (s *LocalFileStorage) SaveFile(fileData io.Reader, filename string, directory string) (string, error) {
	// 保存先ディレクトリのパスを作成
	dirPath := filepath.Join(s.basePath, directory)

	// ディレクトリが存在しない場合は作成
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// ファイルパスを作成
	filePath := filepath.Join(dirPath, filename)

	// ファイルを作成
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// ファイルデータを書き込み
	_, err = io.Copy(file, fileData)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// 相対パスを返す
	return filepath.Join(directory, filename), nil
}

// DeleteFile はファイルを削除します
func (s *LocalFileStorage) DeleteFile(filePath string) error {
	// 絶対パスを作成
	absPath := filepath.Join(s.basePath, filePath)

	// ファイルが存在するか確認
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// ファイルを削除
	if err := os.Remove(absPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFilePath はファイルの絶対パスを取得します
func (s *LocalFileStorage) GetFilePath(filePath string) string {
	return filepath.Join(s.basePath, filePath)
}
