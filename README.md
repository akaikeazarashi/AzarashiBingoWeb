# アザラシビンゴwebアプリ

アザラシに関するビンゴに挑戦できるWebアプリケーションです。

https://bingo.azarashiapp.com/

## 使われている技術

### バックエンド
* **言語**: Go 1.25.3
* **Webフレームワーク**: Gin
* **ORM**: GORM
* **データベース**: MySQL
* **認証**: JWT (golang-jwt/jwt/v5)
* **マイグレーション**: gormigrate
* **環境変数管理**: godotenv

### フロントエンド
* **フレームワーク**: Vue 3
* **ビルドツール**: Vite
* **ルーティング**: Vue Router 4
* **状態管理**: Pinia
* **UIフレームワーク**: Bootstrap 5
* **HTTPクライアント**: Axios
* **その他ライブラリ**:
  * html2canvas (画像キャプチャ)
  * vue-toastification (トースト通知)
  * @kouts/vue-modal (モーダル)

### インフラ
* **CDN**: CloudFront
* **ストレージ**: AWS S3

## プロジェクト構成
※セキュリティのためいくつかコミットしていないファイルがあります

```
AzarashiBingoWeb/
├── app/
│   ├── models/          # データモデル
│   ├── repositories/    # データアクセス層
│   ├── services/        # ビジネスロジック層
│   └── util/            # ユーティリティ
├── bat/                 # デプロイ用バッチファイル (未コミット)
├── bin/                 # バックエンドビルド成果物
├── cmd/
│   └── createUser/      # 管理ユーザー作成コマンド
├── config/              # アプリケーション設定 (未コミット)
├── database/            # データベース初期化・マイグレーション
├── front/               # フロントエンド
│   ├── src/
│   │   ├── components/  # Vueコンポーネント
│   │   ├── views/       # ページビュー
│   │   ├── stores/      # Piniaストア
│   │   └── util/        # ユーティリティ
│   └── dist/            # フロントエンドビルド成果物
├── route/               # ルーティング設定
└── resources/           # 静的リソース
```

## 機能

### ユーザー向け機能
* ビンゴ一覧表示
* ビンゴ詳細表示・プレイ
* ビンゴ結果の送信
* ビンゴ画像のx(Twitter)投稿

### 管理画面機能
* ビンゴアイテムの一覧・詳細表示
* ビンゴアイテムの登録・更新・削除
* ビンゴアイテムのインポート・エクスポート
