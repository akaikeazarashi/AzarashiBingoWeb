package repositories_test

import (
	"AzarashiBingoWeb/app/models"
	"AzarashiBingoWeb/app/repositories"
	"AzarashiBingoWeb/database"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetBingo(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	mockId := uint(1)
	mockName := "ビンゴテスト1"
	mockDescription := "モック説明です"
	mockSize := 3

	rows := sqlmock.NewRows([]string{"ID", "Name", "Description", "Size"}).
		AddRow(mockId, mockName, mockDescription, mockSize)
	mock.ExpectQuery("SELECT \\* FROM `bingos` WHERE id = \\?").
		WillReturnRows(rows)

	mockItemId := uint(1)
	mockItemName := "ビンゴアイテムテスト1"
	mockOrderIndex := 1

	itemRows := sqlmock.NewRows([]string{"ID", "BingoID", "OrderIndex", "Name"}).
		AddRow(mockItemId, mockId, mockOrderIndex, mockItemName)
	mock.ExpectQuery("SELECT \\* FROM `bingo_items` WHERE `bingo_items`.`bingo_id` = \\?").
		WillReturnRows(itemRows)

	result, err := repositories.GetBingo(db, mockId)
	if err != nil {
		t.Fatalf("get error: %s", err)
	}

	if result.ID != mockId {
		t.Errorf("expected ID %d, got %d", mockId, result.ID)
	}
	if result.Name != mockName {
		t.Errorf("expected Name %s, got %s", mockName, result.Name)
	}
	if result.Description != mockDescription {
		t.Errorf("expected Description %s, got %s", mockDescription, result.Description)
	}
	if result.Size != mockSize {
		t.Errorf("expected Size %d, got %d", mockSize, result.Size)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestGetBingoList(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	mockBingos := []struct {
		ID          uint
		Name        string
		Description string
		Size        int
	}{
		{ID: 1, Name: "ビンゴテスト1", Description: "モック説明です", Size: 3},
		{ID: 2, Name: "ビンゴテスト2", Description: "モック説明2です", Size: 5},
		{ID: 3, Name: "ビンゴテスト3", Description: "モック説明3です", Size: 7},
	}

	rows := sqlmock.NewRows([]string{"ID", "Name", "Description", "Size"})
	for _, bingo := range mockBingos {
		rows.AddRow(bingo.ID, bingo.Name, bingo.Description, bingo.Size)
	}
	mock.ExpectQuery("SELECT \\* FROM `bingos` ORDER BY id desc").WillReturnRows(rows)

	// 各ビンゴのbingo_items取得クエリ検証
	itemRows := sqlmock.NewRows([]string{"ID", "BingoID", "OrderIndex", "Name"})
	mock.ExpectQuery("SELECT \\* FROM `bingo_items` WHERE `bingo_items`.`bingo_id` IN \\(\\?,\\?,\\?\\) ORDER BY order_index ASC").WillReturnRows(itemRows)

	result, err := repositories.GetBingoList(db)
	if err != nil {
		t.Fatalf("get list error: %s", err)
	}

	// 取得レコード数検証
	expectedCount := len(mockBingos)
	if len(result) != expectedCount {
		t.Errorf("expected %d bingos, got %d", expectedCount, len(result))
	}

	// 各レコードの内容を検証
	for i, expected := range mockBingos {
		if result[i].ID != expected.ID {
			t.Errorf("expected ID %d, got %d", expected.ID, result[i].ID)
		}
		if result[i].Name != expected.Name {
			t.Errorf("expected Name %s, got %s", expected.Name, result[i].Name)
		}
		if result[i].Size != expected.Size {
			t.Errorf("expected Size %d, got %d", expected.Size, result[i].Size)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestCreateBingo(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	bingo := &models.Bingo{
		Name:        "ビンゴテスト1",
		Description: "ビンゴテスト1の説明です",
		Size:        4,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `bingos`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repositories.CreateBingo(db, bingo)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestUpdatesBingo(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	bingo := &models.Bingo{
		ID:          1,
		Name:        "Update Bingo",
		Description: "updatedesc",
		Size:        5,
	}

	values := map[string]interface{}{"Name": "UpdatedName"}
	items := []models.BingoItem{
		{ID: 1, Name: "item1", OrderIndex: 1},
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `bingos` SET").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM `bingo_items` WHERE bingo_id = \\?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `bingo_items`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repositories.UpdatesBingo(db, bingo, values, items)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestDeleteBingo(t *testing.T) {
	db, mock, err := database.InitDBMock(t)
	if err != nil {
		t.Fatalf("mock init error: %s", err)
	}

	bingo := &models.Bingo{ID: 1}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `bingos` WHERE `bingos`.`id` = \\?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repositories.DeleteBingo(db, bingo)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}
