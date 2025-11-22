package util

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
)

// GetContext - 各種処理用のContext取得
func GetContext() (context.Context, context.CancelFunc) {
	timeout, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	if err != nil {
		log.Println(err)
	}

	context, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	return context, cancel
}
