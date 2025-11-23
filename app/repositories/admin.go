package repositories

import (
	"AzarashiBingoWeb/app/models"

	"gorm.io/gorm"
)

// GetAdminUser - 管理画面ユーザー情報取得
func GetAdminUser(db *gorm.DB, userId string) models.AdminUser {
	adminUser := models.AdminUser{}
	db.Where("user_id = ?", userId).First(&adminUser)

	return adminUser
}
