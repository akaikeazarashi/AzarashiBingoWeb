package services

import (
	"AzarashiBingoWeb/app/models"
	"AzarashiBingoWeb/app/repositories"
	"AzarashiBingoWeb/app/util"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetItemList - ビンゴ一覧を取得
func GetItemList(c *gin.Context, db *gorm.DB) {
	bingoList, err := repositories.GetBingoList(db)

	if err != nil {
		log.Fatalf("Failed GetItemList: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "bingoList": bingoList})
}

// GetItem - ビンゴ詳細取得
func GetItem(c *gin.Context, db *gorm.DB) {
	paramBingoId := c.Param("id")

	bingo, err := repositories.GetBingo(db, paramBingoId)
	if err != nil {
		log.Printf("Failed fetch Bingo record: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}

	if err != nil {
		log.Fatalf("Failed GetItem: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "item": bingo})
}

// SubmitResult - ビンゴ結果の計算
func SubmitResult(c *gin.Context, db *gorm.DB) {
	requestData := models.RequestSubmitResult{}
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		log.Printf("SubmitResult: Invalid Request JSON")
		c.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}

	if requestData.BingoId == 0 {
		log.Printf("SubmitResult: Invalid BingoId")
		c.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}

	// リクエスト情報をマップに直す
	requestMap := make(map[int]bool)
	for _, requestItem := range requestData.BingoItems {
		requestMap[requestItem.Id] = requestItem.IsChecked
	}

	// idに対応するビンゴ情報をDBから取得
	bingo, err := repositories.GetBingo(db, requestData.BingoId)
	if err != nil {
		log.Printf("Failed fetch Bingo record: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}

	// ビンゴアイテムを二次元配列に直す
	bingoSize := bingo.Size
	bingoItemList := make([][]BingoResultItem, bingoSize)
	for i := range bingoItemList {
		bingoItemList[i] = make([]BingoResultItem, bingoSize)
	}

	row := 0
	col := 0
	for _, item := range bingo.Items {
		itemId := item.ID

		// リクエストの達成情報をマージ
		checkedResult, isExist := requestMap[int(itemId)]
		isChecked := false
		if isExist {
			isChecked = checkedResult
		}

		resultItem := BingoResultItem{
			ItemId:    int(itemId),
			IsChecked: isChecked,
		}

		bingoItemList[row][col] = resultItem
		col++

		// colのインクリメント後で判断するのでbingoSizeに-1は不要
		if col == (bingoSize) {
			row++
			col = 0
		}
	}

	rowAchievedIndexList := []int{}
	colAchievedIndexList := []int{}
	diagonalAchievedIndexList := []int{}
	achievedItemIdList := []int{}

	// ビンゴ達成判定(横)
	for rowIndex := 0; rowIndex < bingoSize; rowIndex++ {
		isAchieved := true
		for colIndex := 0; colIndex < bingoSize; colIndex++ {
			if !bingoItemList[rowIndex][colIndex].IsChecked {
				isAchieved = false
				break
			}
		}

		if isAchieved {
			rowAchievedIndexList = append(rowAchievedIndexList, rowIndex)
			for index := 0; index < bingoSize; index++ {
				achievedItemIdList = append(achievedItemIdList, bingoItemList[rowIndex][index].ItemId)
			}
		}
	}

	// ビンゴ達成判定(縦)
	for colIndex := 0; colIndex < bingoSize; colIndex++ {
		isAchieved := true
		for rowIndex := 0; rowIndex < bingoSize; rowIndex++ {
			if !bingoItemList[rowIndex][colIndex].IsChecked {
				isAchieved = false
				break
			}
		}

		if isAchieved {
			colAchievedIndexList = append(colAchievedIndexList, colIndex)
			for index := 0; index < bingoSize; index++ {
				achievedItemIdList = append(achievedItemIdList, bingoItemList[index][colIndex].ItemId)
			}
		}
	}

	// ビンゴ達成判定(左斜め上→右斜め下)
	tempAchievedItemIdList := []int{}
	isAchieved := true
	for diagonalIndex := 0; diagonalIndex < bingoSize; diagonalIndex++ {
		bingoItem := bingoItemList[diagonalIndex][diagonalIndex]
		if !bingoItem.IsChecked {
			isAchieved = false
			break
		}
		tempAchievedItemIdList = append(tempAchievedItemIdList, bingoItem.ItemId)
	}
	if isAchieved {
		diagonalAchievedIndexList = append(diagonalAchievedIndexList, 0)
		achievedItemIdList = append(achievedItemIdList, tempAchievedItemIdList...)
	}

	// ビンゴ達成判定(右斜め上→左斜め下)
	tempAchievedItemIdList = []int{}
	isAchieved = true
	for diagonalRowIndex := bingoSize - 1; 0 <= diagonalRowIndex; diagonalRowIndex-- {
		diagonalColIndex := absInt(diagonalRowIndex - (bingoSize - 1))
		bingoItem := bingoItemList[diagonalRowIndex][diagonalColIndex]
		if !bingoItem.IsChecked {
			isAchieved = false
			break
		}
		tempAchievedItemIdList = append(tempAchievedItemIdList, bingoItem.ItemId)
	}
	if isAchieved {
		diagonalAchievedIndexList = append(diagonalAchievedIndexList, 1)
		achievedItemIdList = append(achievedItemIdList, tempAchievedItemIdList...)
	}

	// achievedItemIdListの重複削除
	slices.Sort(achievedItemIdList)
	achievedItemIdList = slices.Compact(achievedItemIdList)

	achievedCount := len(rowAchievedIndexList) + len(colAchievedIndexList) + len(diagonalAchievedIndexList)
	isAllAchieved := len(bingo.Items) <= len(achievedItemIdList)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"bingoResult": gin.H{
			"rowAchievedIndexList":           rowAchievedIndexList,
			"colAchievedIndexList":           colAchievedIndexList,
			"diagonalAchievedgonalIndexList": diagonalAchievedIndexList,
			"achievedCount":                  achievedCount,
			"achievedItemIdList":             achievedItemIdList,
			"isAllAchieved":                  isAllAchieved,
		},
	})
}

// UploadBingoImage - キャプチャしたビンゴ画像をS3へアップロード
func UploadBingoImage(c *gin.Context, db *gorm.DB) {
	// ビンゴ情報取得
	paramBingoId := c.PostForm("bingoId")

	bingo, err := repositories.GetBingo(db, paramBingoId)
	if err != nil {
		log.Printf("Failed fetch Bingo record: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}

	// アップロードファイル取得
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	file := form.File["file"][0]

	imageSrc, err := file.Open()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer imageSrc.Close()

	// 処理用Context生成
	context, cancel := util.GetContext()
	defer cancel()

	// AWS設定ロード
	awsConfig, err := util.LoadAWSConfig(context)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "aws config load failed"})
		return
	}

	env := os.Getenv("ENV")
	cloudfrontDomainName := os.Getenv("AWS_CLOUDFRONT_DOMAIN_NAME")
	dirName := time.Now().Format("20060102_030405.000")
	htmlKey := fmt.Sprintf("%s/%s/share.html", env, dirName)
	htmlUrl := fmt.Sprintf("https://%s/%s", cloudfrontDomainName, htmlKey)
	imageKey := fmt.Sprintf("%s/%s/image%s", env, dirName, filepath.Ext(file.Filename))
	imageUrl := fmt.Sprintf("https://%s/%s", cloudfrontDomainName, imageKey)

	// 画像をS3へアップロード
	s3Client := s3.NewFromConfig(awsConfig)
	bucket := os.Getenv("AWS_S3_BUCKET")

	_, err = s3Client.PutObject(context, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(imageKey),
		Body:        imageSrc,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "aws s3 put image failed"})
		return
	}

	// OGP用HTMLを整形して取得
	htmlText, err := getShareHtmlText(bingo, htmlUrl, imageUrl)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "share html error"})
		return
	}

	// OGP用HTMLをS3へアップロード
	_, err = s3Client.PutObject(context, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(htmlKey),
		Body:        bytes.NewReader([]byte(htmlText)),
		ContentType: aws.String("text/html; charset=utf-8"),
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "aws s3 put html failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "url": htmlUrl})
}

// getShareHtmlText - SNS投稿用HTMLファイルのテキストを取得する
func getShareHtmlText(bingo *models.Bingo, htmlUrl string, imageUrl string) (string, error) {
	htmlBytes, err := os.ReadFile("resources/html/share.html")
	if err != nil {
		return "", err
	}
	html := string(htmlBytes)

	appName := os.Getenv("APP_NAME")

	html = strings.ReplaceAll(html, "%title%", fmt.Sprintf("【%s】%s", appName, bingo.Name))
	html = strings.ReplaceAll(html, "%image_url%", imageUrl)
	html = strings.ReplaceAll(html, "%page_url%", htmlUrl)

	// Descriptionは改行が含まれる可能性があるので除去
	reg := regexp.MustCompile(`\r?\n`)
	description := reg.ReplaceAllString(bingo.Description, "")
	html = strings.ReplaceAll(html, "%description%", description)

	siteDomain := os.Getenv("SITE_DOMAIN")
	redirectUrl := fmt.Sprintf("https://%s/bingo/%d", siteDomain, bingo.ID)
	html = strings.ReplaceAll(html, "%redirect_url%", redirectUrl)

	return html, nil
}

// absInt - 整数の絶対値を返す
func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// BingoResultItem - ビンゴ結果返却用のアイテム構造
type BingoResultItem struct {
	ItemId    int
	IsChecked bool
}
