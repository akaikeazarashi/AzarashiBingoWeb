package models

// RequestSubmitResult - ビンゴデータ送信リクエスト
type RequestSubmitResult struct {
	BingoId    int                           `json:"id"`
	BingoItems []RequestSubmitResultItemInfo `json:"items"`
}

// RequestSubmitResultItemInfo - ビンゴデータ送信のビンゴアイテム情報
type RequestSubmitResultItemInfo struct {
	Id        int  `json:"id"`
	IsChecked bool `json:"isChecked"`
}

// RequestSingin - 管理画面ログインリクエスト
type RequestAdminSingin struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

// RequestAdminItemPut - 管理画面アイテム登録・更新
type RequestAdminItemPut struct {
	BingoId     int                         `json:"bingoId"`
	BingoName   string                      `json:"name"`
	Description string                      `json:"description"`
	Size        int                         `json:"size"`
	BingoItems  []RequestAdminBingoItemInfo `json:"items"`
}

// RequestAdminItemImport - 管理画面アイテムインポート
type RequestAdminItemImport struct {
	BingoName   string                      `json:"name"`
	Description string                      `json:"description"`
	Size        int                         `json:"size"`
	BingoItems  []RequestAdminBingoItemInfo `json:"items"`
}

// ItemInfo - ビンゴアイテム情報
type RequestAdminBingoItemInfo struct {
	Name       string `json:"name"`
	OrderIndex int    `json:"orderIndex"`
}
