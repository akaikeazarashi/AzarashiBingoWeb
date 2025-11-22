package repositories

import (
	"AzarashiBingoWeb/app/models"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// GetBingo - idに対応するビンゴ情報を取得する
func GetBingo(db *gorm.DB, id any) (*models.Bingo, error) {
	var bingoId int

	switch idValue := id.(type) {
	case int:
		bingoId = idValue
	case uint:
		bingoId = int(idValue)
	case string:
		idInt, err := strconv.Atoi(idValue)
		if err != nil {
			return nil, fmt.Errorf("invalid string for int: %v", err)
		}
		bingoId = idInt
	default:
		return nil, fmt.Errorf("unsupported type: %T", id)
	}

	var bingo models.Bingo
	result := db.Preload("Items", models.OrderbyOrderIndex).Where("id = ?", bingoId).Order("id desc").First(&bingo)

	return &bingo, result.Error
}

// GetBingoList - 全ビンゴ情報を取得する
func GetBingoList(db *gorm.DB) ([]models.Bingo, error) {
	var bingoList []models.Bingo
	result := db.Preload("Items", models.OrderbyOrderIndex).Order("id desc").Find(&bingoList)

	return bingoList, result.Error
}

// CreateBingo - ビンゴデータ作成
func CreateBingo(db *gorm.DB, bingo *models.Bingo) error {
	result := db.Create(&bingo)
	return result.Error
}

// UpdatesBingo - ビンゴデータ更新
func UpdatesBingo(db *gorm.DB, bingo *models.Bingo, values map[string]interface{}, items []models.BingoItem) error {
	tx := db.Begin()

	// ビンゴ情報更新
	// (Saveだとcreated_at等の更新も走るのでUpdatesで行う)
	updateResult := tx.Model(&bingo).Updates(values)
	if updateResult.Error != nil {
		tx.Rollback()
		return updateResult.Error
	}

	// アイテムをすべて削除し追加し直す
	deleteResult := tx.Delete(&models.BingoItem{}, "bingo_id = ?", bingo.ID)
	if deleteResult.Error != nil {
		tx.Rollback()
		return deleteResult.Error
	}
	for _, item := range items {
		item.BingoID = bingo.ID
		itemCreateResult := tx.Create(&item)

		if itemCreateResult.Error != nil {
			tx.Rollback()
			return itemCreateResult.Error
		}
	}

	return tx.Commit().Error
}

// DeleteBingo - ビンゴデータ削除
func DeleteBingo(db *gorm.DB, bingo *models.Bingo) error {
	result := db.Delete(&bingo)
	return result.Error
}
