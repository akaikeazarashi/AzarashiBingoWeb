package models

import (
	"time"

	"gorm.io/gorm"
)

// Bingo - ビンゴデータ
type Bingo struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Size        int         `json:"size"`
	Items       []BingoItem `gorm:"foreignKey:BingoID;constraint:OnDelete:CASCADE;" json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// BingoItem - ビンゴの1マスあたりのアイテム
// Bingoに対して1-nレコード存在する
type BingoItem struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	BingoID    uint   `gorm:"index" json:"bingoId"`
	Name       string `json:"name"`
	OrderIndex int    `json:"orderIndex"`
}

// OrderbyOrderIndex - BingoItemのOrder指定用のスコープ関数
func OrderbyOrderIndex(db *gorm.DB) *gorm.DB {
	return db.Order("order_index ASC")
}
