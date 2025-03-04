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
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// パスパラメータから温泉IDを取得
	onsenID := ctx.Param("onsen_id")
	if onsenID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ONSEN_ID",
				"message": "温泉IDが指定されていません",
			},
		})
		return
	}

	// ファイルを取得
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_FILE",
				"message": "画像ファイルが無効です: " + err.Error(),
			},
		})
		return
	}
	defer file.Close()

	// 説明文を取得
	description := ctx.PostForm("description")

	// 入力データを作成
	input := port.UploadImageInput{
		OnsenID:     onsenID,
		UserID:      userID.(string),
		File:        file,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Description: description,
	}

	// ユースケースを呼び出し
	image, err := c.onsenImageUseCase.UploadImage(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPLOAD_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":          image.ID,
			"onsen_id":    image.OnsenID,
			"url":         image.URL,
			"description": image.Description,
			"created_at":  image.CreatedAt,
		},
		"message": "画像のアップロードに成功しました",
	})
}

// GetImagesByOnsenID は温泉IDに紐づく画像を取得します
func (c *OnsenImageController) GetImagesByOnsenID(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// パスパラメータから温泉IDを取得
	onsenID := ctx.Param("onsen_id")
	if onsenID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ONSEN_ID",
				"message": "温泉IDが指定されていません",
			},
		})
		return
	}

	// 入力データを作成
	input := port.GetImagesByOnsenIDInput{
		OnsenID: onsenID,
		UserID:  userID.(string),
	}

	// ユースケースを呼び出し
	images, err := c.onsenImageUseCase.GetImagesByOnsenID(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "FETCH_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	// レスポンスデータを作成
	var imagesData []gin.H
	for _, img := range images {
		imagesData = append(imagesData, gin.H{
			"id":          img.ID,
			"onsen_id":    img.OnsenID,
			"url":         img.URL,
			"description": img.Description,
			"created_at":  img.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"images": imagesData,
		},
		"message": "画像の取得に成功しました",
	})
}

// DeleteImage は温泉画像を削除します
func (c *OnsenImageController) DeleteImage(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "認証されていません",
			},
		})
		return
	}

	// パスパラメータから画像IDを取得
	imageID := ctx.Param("image_id")
	if imageID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_IMAGE_ID",
				"message": "画像IDが指定されていません",
			},
		})
		return
	}

	// 入力データを作成
	input := port.DeleteImageInput{
		ImageID: imageID,
		UserID:  userID.(string),
	}

	// ユースケースを呼び出し
	err := c.onsenImageUseCase.DeleteImage(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DELETE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "画像の削除に成功しました",
	})
}
