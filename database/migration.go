package database

import (
	"AzarashiBingoWeb/app/models"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate - DBのマイグレーション実行を試す
func Migrate(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())
	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Println("Migration completed successfully!")
}

// getMigrations - マイグレーションリストを取得する
func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "2025102401_create_bingos_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Bingo{}, &models.BingoItem{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Bingo{}, &models.BingoItem{})
			},
		},
		{
			ID: "2025031301_create_admin_users_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.AdminUser{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.AdminUser{})
			},
		},
	}
}
