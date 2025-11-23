package main

import (
	"AzarashiBingoWeb/app/models"
	"AzarashiBingoWeb/database"
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

// go run .\cmd\createUser\main.go createUser --userid=user_id --password=password
func main() {
	var userid string
	var password string

	var rootCmd = &cobra.Command{
		Use:   "createUser",
		Short: "Create a new admin user",
		Run: func(cmd *cobra.Command, args []string) {
			if userid == "" || password == "" {
				fmt.Println("ユーザーidかパスワードが未入力")
				return
			}

			// 環境変数の初期化
			err := godotenv.Load()
			if err != nil {
				panic("Failed to load env file")
			}

			// DB初期化
			db := database.InitDB()

			adminUser := models.AdminUser{}
			db.Where("user_id = ?", userid).Find(&adminUser)

			// 同一idのユーザーがいないかチェック
			if adminUser.UserId != "" {
				err := errors.New("同一名のuserIdが既に登録されています")
				fmt.Println(err)
				return
			}

			// パスワードをハッシュ化
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println("パスワード暗号化エラー: ", err)
				return
			}

			// ユーザー情報作成
			newUser := models.AdminUser{
				UserId:   userid,
				Password: string(hashedPassword),
			}

			db.Create(&newUser)
			fmt.Println("Admin user created successfully!: " + userid)
		},
	}

	rootCmd.Flags().StringVar(&userid, "userid", "", "admin user_id")
	rootCmd.Flags().StringVar(&password, "password", "", "admin password")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
