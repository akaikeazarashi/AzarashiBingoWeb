package services_test

import (
	"AzarashiBingoWeb/app/models"
	"AzarashiBingoWeb/app/services"
	"AzarashiBingoWeb/database"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

type MockBingo struct {
	ID          uint
	Name        string
	Description string
	Size        int
}

type MockBingoItem struct {
	ID         uint
	BingoId    uint
	Name       string
	OrderIndex int
}

func TestGetItemList(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	// モックデータ準備
	mockBingos := []MockBingo{
		{ID: 1, Name: "ビンゴテスト1", Description: "説明1", Size: 3},
		{ID: 2, Name: "ビンゴテスト2", Description: "説明2", Size: 5},
	}
	rows := sqlmock.NewRows([]string{"ID", "Name", "Description", "Size"})
	for _, bingo := range mockBingos {
		rows.AddRow(bingo.ID, bingo.Name, bingo.Description, bingo.Size)
	}
	mock.ExpectQuery(".+").WillReturnRows(rows)

	mockBingoItems := []MockBingoItem{
		{ID: 1, BingoId: 1, Name: "ビンゴアイテムテスト1-1", OrderIndex: 1},
		{ID: 2, BingoId: 2, Name: "ビンゴアイテムテスト2-1", OrderIndex: 1},
		{ID: 3, BingoId: 2, Name: "ビンゴアイテムテスト2-2", OrderIndex: 2},
	}
	itemRows := sqlmock.NewRows([]string{"ID", "BingoID", "Name", "OrderIndex"})
	for _, bingoItem := range mockBingoItems {
		rows.AddRow(bingoItem.ID, bingoItem.BingoId, bingoItem.Name, bingoItem.OrderIndex)
	}
	mock.ExpectQuery(".+").WillReturnRows(itemRows)

	// リクエスト実行
	c, r := setupTestContext()
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	services.GetItemList(c, db)

	if r.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, r.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(r.Body.Bytes(), &response)
	if response["result"] != true {
		t.Errorf("expected result true, got %v", response["result"])
	}
}

func TestGetItem(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	// モックデータ準備
	mockId := uint(1)
	mockName := "ビンゴテスト1"
	mockDescription := "説明です"
	mockSize := 3

	rows := sqlmock.NewRows([]string{"ID", "Name", "Description", "Size"}).
		AddRow(mockId, mockName, mockDescription, mockSize)
	mock.ExpectQuery(".+").WillReturnRows(rows)

	mockBingoItems := []MockBingoItem{
		{ID: 1, BingoId: 1, Name: "ビンゴアイテムテスト1-1", OrderIndex: 1},
		{ID: 2, BingoId: 1, Name: "ビンゴアイテムテスト1-2", OrderIndex: 2},
	}
	itemRows := sqlmock.NewRows([]string{"ID", "BingoID", "Name", "OrderIndex"})
	for _, bingoItem := range mockBingoItems {
		itemRows.AddRow(bingoItem.ID, bingoItem.BingoId, bingoItem.Name, bingoItem.OrderIndex)
	}
	mock.ExpectQuery(".+").WillReturnRows(itemRows)

	// リクエスト実行
	c, r := setupTestContext()
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	services.GetItem(c, db)

	if r.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, r.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(r.Body.Bytes(), &response)
	if response["result"] != true {
		t.Errorf("expected result true, got %v", response["result"])
	}
}

func TestSubmitResult_Success(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	mockBingoId := uint(1)
	mockSize := 3

	// モックデータ準備
	bingoRows := sqlmock.NewRows([]string{"ID", "Name", "Description", "Size"}).
		AddRow(mockBingoId, "テストビンゴ", "説明", mockSize)
	mock.ExpectQuery(".+").WillReturnRows(bingoRows)

	mockBingoItems := []MockBingoItem{}
	for i := 1; i <= mockSize*mockSize; i++ {
		mockBingoItems = append(mockBingoItems, MockBingoItem{
			ID: uint(i), BingoId: 1, Name: fmt.Sprintf("ビンゴアイテムテスト%d", i), OrderIndex: i - 1,
		})
	}
	itemRows := sqlmock.NewRows([]string{"ID", "BingoID", "Name", "OrderIndex"})
	for _, bingoItem := range mockBingoItems {
		itemRows.AddRow(bingoItem.ID, bingoItem.BingoId, bingoItem.Name, bingoItem.OrderIndex)
	}
	mock.ExpectQuery(".+").WillReturnRows(itemRows)

	// リクエストデータ(1行目と2列目をすべてチェック)
	requestData := models.RequestSubmitResult{
		BingoId: int(mockBingoId),
		BingoItems: []models.RequestSubmitResultItemInfo{
			{Id: 1, IsChecked: true},
			{Id: 2, IsChecked: true},
			{Id: 3, IsChecked: true},
			{Id: 5, IsChecked: true},
			{Id: 8, IsChecked: true},
		},
	}
	jsonBody, _ := json.Marshal(requestData)

	// リクエスト実行
	c, r := setupTestContext()
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	services.SubmitResult(c, db)

	if r.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, r.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(r.Body.Bytes(), &response)
	if response["result"] != true {
		t.Errorf("expected result true, got %v", response["result"])
	}

	bingoResult, ok := response["bingoResult"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected bingoResult in response")
	}

	// ビンゴ達成状態の検証(横1縦1)
	rowAchievedList, ok := bingoResult["rowAchievedIndexList"].([]interface{})
	if !ok {
		t.Fatalf("expected rowAchievedIndexList in bingoResult")
	}
	if len(rowAchievedList) != 1 {
		t.Errorf("expected 1 row achieved, actual: %d", len(rowAchievedList))
	}
	colAchievedList, ok := bingoResult["colAchievedIndexList"].([]interface{})
	if !ok {
		t.Fatalf("expected colAchievedIndexList in bingoResult")
	}
	if len(colAchievedList) != 1 {
		t.Errorf("expected 1 col achieved, actual: %d", len(colAchievedList))
	}
}

// setupTestContext - テスト用のコンテキスト取得
func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	return c, w
}
