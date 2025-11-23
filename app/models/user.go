package models

// AdminUser - 管理画面用ユーザー
type AdminUser struct {
	ID       uint   `gorm:"primaryKey"`
	UserId   string `gorm:"unique"`
	Password string
}
