package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// LocalFileStorage はローカルファイルシステムを使用したストレージ実装です
type LocalFileStorage struct {
	uploadDir string
}

// NewLocalFileStorage は新しいLocalFileStorageインスタンスを作成します
func NewLocalFileStorage(uploadDir string) (*LocalFileStorage, error) {
	// ベースディレクトリが存在しない場合は作成
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, fmt.Errorf("ストレージディレクトリの作成に失敗しました: %w", err)
		}
	}

	return &LocalFileStorage{
		uploadDir: uploadDir,
	}, nil
}

// Upload はファイルをアップロードします
func (s *LocalFileStorage) Upload(ctx context.Context, file io.Reader, fileName, contentType string) (string, error) {
	// ユニークなファイル名を生成
	ext := filepath.Ext(fileName)
	uniqueFileName := uuid.New().String() + ext

	// ファイルパスを作成
	filePath := filepath.Join(s.uploadDir, uniqueFileName)

	// ファイルを作成
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	defer dst.Close()

	// ファイルデータを書き込み
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
	}

	// ファイルURLを返す
	return "/uploads/" + uniqueFileName, nil
}

// Delete はファイルを削除します
func (s *LocalFileStorage) Delete(ctx context.Context, fileURL string) error {
	// ファイルURLからパスを抽出
	filePath := s.extractPathFromURL(fileURL)
	if filePath == "" {
		return fmt.Errorf("無効なファイルURLです: %s", fileURL)
	}

	// ファイルが存在するか確認
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("ファイルが存在しません: %s", filePath)
	}

	// ファイルを削除
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("ファイルの削除に失敗しました: %w", err)
	}

	return nil
}

// extractPathFromURL はURLからファイルパスを抽出します
func (s *LocalFileStorage) extractPathFromURL(fileURL string) string {
	// URLからファイル名を抽出
	parts := strings.Split(fileURL, "/")
	if len(parts) == 0 {
		return ""
	}
	fileName := parts[len(parts)-1]

	// ファイルパスを構築
	return filepath.Join(s.uploadDir, fileName)
}

// SaveFile はファイルを保存します（互換性のため残しています）
func (s *LocalFileStorage) SaveFile(fileData io.Reader, filename string, directory string) (string, error) {
	// 保存先ディレクトリのパスを作成
	dirPath := filepath.Join(s.uploadDir, directory)

	// ディレクトリが存在しない場合は作成
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
		}
	}

	// ファイルパスを作成
	filePath := filepath.Join(dirPath, filename)

	// ファイルを作成
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	defer file.Close()

	// ファイルデータを書き込み
	_, err = io.Copy(file, fileData)
	if err != nil {
		return "", fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
	}

	// 相対パスを返す
	return filepath.Join(directory, filename), nil
}

// DeleteFile はファイルを削除します（互換性のため残しています）
func (s *LocalFileStorage) DeleteFile(filePath string) error {
	// 絶対パスを作成
	absPath := filepath.Join(s.uploadDir, filePath)

	// ファイルが存在するか確認
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("ファイルが存在しません: %s", filePath)
	}

	// ファイルを削除
	if err := os.Remove(absPath); err != nil {
		return fmt.Errorf("ファイルの削除に失敗しました: %w", err)
	}

	return nil
}

// GetFilePath はファイルの絶対パスを取得します
func (s *LocalFileStorage) GetFilePath(filePath string) string {
	return filepath.Join(s.uploadDir, filePath)
}
