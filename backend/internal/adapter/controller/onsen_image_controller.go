package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yuroku/internal/usecase/port"
)

// OnsenImageController は温泉画像関連のコントローラーです
type OnsenImageController struct {
	onsenImageUseCase port.OnsenImageInputPort
}

// NewOnsenImageController は新しい温泉画像コントローラーを作成します
func NewOnsenImageController(onsenImageUseCase port.OnsenImageInputPort) *OnsenImageController {
	return &OnsenImageController{
		onsenImageUseCase: onsenImageUseCase,
	}
}

// UploadImage は温泉画像をアップロードします
func (c *OnsenImageController) UploadImage(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// パスパラメータから温泉IDを取得
	onsenID, ok := ValidatePathParam(ctx, "onsen_id", "温泉IDが指定されていません")
	if !ok {
		return
	}

	// ファイルを取得
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		RespondWithError(ctx, http.StatusBadRequest, "INVALID_FILE", "画像ファイルが無効です: "+err.Error())
		return
	}
	defer file.Close()

	// 説明文を取得
	description := ctx.PostForm("description")

	// 入力データを作成
	input := port.UploadImageInput{
		OnsenID:     onsenID,
		UserID:      userID,
		File:        file,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Description: description,
	}

	// ユースケースを呼び出し
	image, err := c.onsenImageUseCase.UploadImage(ctx.Request.Context(), input)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusCreated, gin.H{
		"id":          image.ID,
		"onsen_id":    image.OnsenID,
		"url":         image.URL,
		"description": image.Description,
		"created_at":  image.CreatedAt,
	}, "画像のアップロードに成功しました")
}

// GetImagesByOnsenID は温泉IDに紐づく画像を取得します
func (c *OnsenImageController) GetImagesByOnsenID(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// パスパラメータから温泉IDを取得
	onsenID, ok := ValidatePathParam(ctx, "onsen_id", "温泉IDが指定されていません")
	if !ok {
		return
	}

	// 入力データを作成
	input := port.GetImagesByOnsenIDInput{
		OnsenID: onsenID,
		UserID:  userID,
	}

	// ユースケースを呼び出し
	images, err := c.onsenImageUseCase.GetImagesByOnsenID(ctx.Request.Context(), input)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	// レスポンスデータを作成
	imagesData := make([]gin.H, 0, len(images))
	for _, img := range images {
		imagesData = append(imagesData, gin.H{
			"id":          img.ID,
			"onsen_id":    img.OnsenID,
			"url":         img.URL,
			"description": img.Description,
			"created_at":  img.CreatedAt,
		})
	}

	RespondWithSuccess(ctx, http.StatusOK, gin.H{
		"images": imagesData,
	}, "画像の取得に成功しました")
}

// DeleteImage は温泉画像を削除します
func (c *OnsenImageController) DeleteImage(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, ok := GetUserID(ctx)
	if !ok {
		return
	}

	// パスパラメータから画像IDを取得
	imageID, ok := ValidatePathParam(ctx, "image_id", "画像IDが指定されていません")
	if !ok {
		return
	}

	// 入力データを作成
	input := port.DeleteImageInput{
		ImageID: imageID,
		UserID:  userID,
	}

	// ユースケースを呼び出し
	err := c.onsenImageUseCase.DeleteImage(ctx.Request.Context(), input)
	if err != nil {
		RespondWithAppError(ctx, err)
		return
	}

	RespondWithSuccess(ctx, http.StatusOK, nil, "画像の削除に成功しました")
}
