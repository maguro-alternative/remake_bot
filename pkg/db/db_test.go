package db

import (
	"context"
	"fmt"
	"os"
	"time"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func TestNewDB(t *testing.T) {
	t.Run("DB接続が正常に行われること", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		dbUri := fmt.Sprintf("%s://%s:%s/%s?user=%s&password=%s&sslmode=disable", os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
		_, cleanup, err := NewDB(ctx, "postgres", dbUri)
		assert.NoError(t, err)
		defer cleanup()
	})

	t.Run("DB接続が失敗すること", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		_, cleanup, err := NewDB(ctx, "postgres", "hoge")
		assert.Error(t, err)
		defer cleanup()
	})

	t.Run("DB接続がタイムアウトすること", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		dbUri := fmt.Sprintf("%s://%s:%s/%s?user=%s&password=%s&sslmode=disable", os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
		_, cleanup, err := NewDB(ctx, "postgres", dbUri)
		assert.Error(t, err)
		defer cleanup()
	})
}
