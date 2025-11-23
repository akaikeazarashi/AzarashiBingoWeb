package services

import (
	"AzarashiBingoWeb/app/models"
	"AzarashiBingoWeb/app/repositories"
	"AzarashiBingoWeb/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminSignIn - 管理画面のログイン処理を実行しトークンを返す
func AdminSignIn(c *gin.Context, db *gorm.DB) {
	// 許可IPチェック
	clientIP := c.ClientIP()
	isValidIp := false
	for _, ip := range config.AdminAllowedIPs {
		if clientIP == ip {
			isValidIp = true
			break
		}
	}
	if !isValidIp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth error"})
		c.Abort()
	}

	requestAdminSingin := models.RequestAdminSingin{}
	err := c.ShouldBindJSON(&requestAdminSingin)
	if err != nil {
		abort(c, fmt.Sprintf("Invalid Request JSON: %s", err.Error()))
		return
	}

	adminUser := repositories.GetAdminUser(db, requestAdminSingin.UserId)
	db.Where("user_id = ?", requestAdminSingin.UserId).First(&adminUser)

	if adminUser.UserId == "" {
		abort(c, "SignIn: user not found")
		return
	}

	// パスワード検証
	hashedPassword := adminUser.Password
	inputPassword := requestAdminSingin.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		abort(c, "SignIn: invalid user "+err.Error())
		return
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = requestAdminSingin.UserId
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["iat"] = time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 管理画面シークレットキー検証
	secretKey := os.Getenv("ADMIN_SECRET_KEY")
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		abort(c, "SignIn: sign failed "+err.Error())
		return
	}

	// トークンをcookieに設定
	domain := os.Getenv("SITE_DOMAIN")
	c.SetCookie("token", signedToken, 3600, "/", domain, true, true)

	c.JSON(http.StatusOK, gin.H{"result": true})
}

// AdminItemList - ビンゴアイテム一覧を取得
func AdminItemList(c *gin.Context, db *gorm.DB) {
	var bingoList []models.Bingo

	bingoList, err := repositories.GetBingoList(db)

	if err != nil {
		log.Fatalf("Failed AdminItemList: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "bingoList": bingoList})
}

// AdminItemDetail - ビンゴアイテム詳細取得
func AdminItemDetail(c *gin.Context, db *gorm.DB) {
	bingoId := c.Param("id")
	bingo, err := repositories.GetBingo(db, bingoId)

	if err != nil {
		log.Fatalf("Failed AdminItemDetail: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "item": bingo})
}

// AdminItemPut - ビンゴアイテム登録・更新
func AdminPutItem(c *gin.Context, db *gorm.DB) {
	requestAdminPut := models.RequestAdminItemPut{}
	err := c.ShouldBindJSON(&requestAdminPut)
	if err != nil {
		abort(c, "AdminItemPut: Invalid Request JSON")
		return
	}

	// BingoItemをモデルの方に合わせて変換
	var items []models.BingoItem
	for _, requestItem := range requestAdminPut.BingoItems {
		items = append(items, models.BingoItem{
			Name:       requestItem.Name,
			OrderIndex: requestItem.OrderIndex,
		})
	}

	// idがポストされていたら更新にする
	editBingoId := requestAdminPut.BingoId
	isEdit := editBingoId != 0
	if isEdit {
		dbBingo, err := repositories.GetBingo(db, editBingoId)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err})
			return
		}

		// ビンゴ情報更新
		// (Saveだとcreated_at等の更新も走るのでUpdatesで行う)
		err = repositories.UpdatesBingo(db, dbBingo, map[string]interface{}{
			"name":        requestAdminPut.BingoName,
			"description": requestAdminPut.Description,
			"size":        requestAdminPut.Size,
		}, items)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err})
			return
		}
	} else {
		bingo := models.Bingo{
			Name:        requestAdminPut.BingoName,
			Description: requestAdminPut.Description,
			Size:        requestAdminPut.Size,
			Items:       items,
		}

		// ビンゴ情報保存
		err := repositories.CreateBingo(db, &bingo)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

// AdminDeleteItem - ビンゴアイテム削除
func AdminDeleteItem(c *gin.Context, db *gorm.DB) {
	bingoId := c.Param("id")

	bingo, err := repositories.GetBingo(db, bingoId)
	if err != nil {
		abort(c, fmt.Sprintf("Failed AdminDelete: %v", err))
		return
	}

	err = repositories.DeleteBingo(db, bingo)
	if err != nil {
		abort(c, fmt.Sprintf("Failed AdminDelete: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

// AdminImportItem - ビンゴアイテムのJSONをインポート
func AdminImportItem(c *gin.Context, db *gorm.DB) {
	requestAdminImport := models.RequestAdminItemImport{}
	err := c.ShouldBindJSON(&requestAdminImport)
	if err != nil {
		abort(c, "Invalid Request JSON")
		return
	}

	// BingoItemをモデルの方に合わせて変換
	var items []models.BingoItem
	for _, requestItem := range requestAdminImport.BingoItems {
		items = append(items, models.BingoItem{
			Name:       requestItem.Name,
			OrderIndex: requestItem.OrderIndex,
		})
	}

	bingo := models.Bingo{
		Name:        requestAdminImport.BingoName,
		Description: requestAdminImport.Description,
		Size:        requestAdminImport.Size,
		Items:       items,
	}

	// ビンゴ情報保存
	err = repositories.CreateBingo(db, &bingo)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

// AdminExportItem - ビンゴアイテムのJSONをエクスポート
func AdminExportItem(c *gin.Context, db *gorm.DB) {
	bingoId := c.Param("id")

	bingo, err := repositories.GetBingo(db, bingoId)
	if err != nil {
		log.Printf("Failed export Bingo item: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}

	jsonBytes, err := json.MarshalIndent(bingo, "", "  ")
	if err != nil {
		log.Printf("Failed export Bingo item: %v", err)
		c.String(http.StatusInternalServerError, "JSON変換エラー")
		return
	}

	// ファイルとしてダウンロードさせるためのヘッダー設定
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=bingo_%s.json", bingoId))
	c.Data(http.StatusOK, "application/json", jsonBytes)
}

// abort - ログ出力後にAbort()し処理中断
func abort(c *gin.Context, message string) {
	fmt.Printf("%v\n", message)
	c.JSON(http.StatusUnauthorized, gin.H{"result": false})
	c.Abort()
}
