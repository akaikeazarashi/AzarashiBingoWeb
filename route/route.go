package route

import (
	"AzarashiBingoWeb/app/services"
	"AzarashiBingoWeb/config"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func SetRoutes(r *gin.Engine, db *gorm.DB) {
	// Vueのビルドファイルのパス指定
	r.Static("/static", config.FrontDistPath)
	r.Static("/images", path.Join(config.FrontDistPath, "images"))

	// トップページのルーティング
	r.GET("/", func(c *gin.Context) {
		routeTop(c)
	})

	// APIのルーティング
	api := r.Group("/api")
	{
		// ビンゴ一覧表示
		api.GET("/items", func(c *gin.Context) {
			services.GetItemList(c, db)
		})

		// ビンゴ詳細取得
		api.GET("/item/:id", func(c *gin.Context) {
			services.GetItem(c, db)
		})

		// ビンゴ結果送信
		api.POST("/submit", func(c *gin.Context) {
			services.SubmitResult(c, db)
		})

		// ビンゴ画像アップロード
		api.POST("/bingo/upload", func(c *gin.Context) {
			services.UploadBingoImage(c, db)
		})
	}

	// 管理画面ログイン
	r.POST("/api/admin/login", func(c *gin.Context) {
		services.AdminSignIn(c, db)
	})

	// 管理画面用の認証が必要なルーティング
	admin := r.Group("/api/admin")
	admin.Use(adminIpAuthMiddleware(), authMiddleware())
	{
		// アイテム一覧
		admin.GET("/list", func(c *gin.Context) {
			services.AdminItemList(c, db)
		})
		// アイテム詳細
		admin.GET("/detail/:id", func(c *gin.Context) {
			services.AdminItemDetail(c, db)
		})
		// アイテム登録・更新
		admin.PUT("/put", func(c *gin.Context) {
			services.AdminPutItem(c, db)
		})
		// アイテム削除
		admin.DELETE("/delete/:id", func(c *gin.Context) {
			services.AdminDeleteItem(c, db)
		})
		// アイテムインポート
		admin.PUT("/import", func(c *gin.Context) {
			services.AdminImportItem(c, db)
		})
		// アイテムエクスポート
		admin.GET("/export/:id", func(c *gin.Context) {
			services.AdminExportItem(c, db)
		})
	}

	// その他すべてのルートをVueのindex.htmlにリダイレクト
	r.NoRoute(func(c *gin.Context) {
		routeTop(c)
	})
}

func routeTop(c *gin.Context) {
	c.File(config.FrontDistPath + "index.html")
}

// authMiddleware - 管理画面用のログイン認証
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "ログインが必要です",
			})
			return
		}

		secretKey := os.Getenv("ADMIN_SECRET_KEY")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			fmt.Println(err.Error())
			fmt.Println(token)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "無効なトークンです",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "トークンの形式が不正です",
			})
		}

		// コンテキストにユーザー情報を保存
		c.Set("user_id", claims["user_id"].(string))

		c.Next()
	}
}

// adminIpAuthMiddleware - 管理画面用の許可IPチェック
func adminIpAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		for _, ip := range config.AdminAllowedIPs {
			if clientIP == ip {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth error"})
		c.Abort()
	}
}
