package main

import (
	"AzarashiBingoWeb/database"
	"AzarashiBingoWeb/route"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 環境変数の初期化
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}

	// DB初期化
	db := database.InitDB()

	// マイグレーションチェック
	database.Migrate(db)

	r := gin.Default()

	// ルーティング設定
	route.SetRoutes(r, db)

	// サーバー起動
	url := os.Getenv("PROXY_DOMAIN")
	port := os.Getenv("PORT")
	r.Run(url + ":" + port)
}
