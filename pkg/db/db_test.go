package db

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
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

func TestDB_Bytea(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbUri := fmt.Sprintf("%s://%s:%s/%s?user=%s&password=%s&sslmode=disable", os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	db, cleanup, err := NewDB(ctx, "postgres", dbUri)
	assert.NoError(t, err)
	tx, err := db.BeginTxx(ctx, nil)
	assert.NoError(t, err)
	defer cleanup()
	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "CREATE TABLE bytea_table (bytea_column BYTEA)")
	assert.NoError(t, err)

	t.Run("Bytea型のデータが正常に登録できること", func(t *testing.T) {
		_, err := tx.ExecContext(ctx, "INSERT INTO bytea_table (bytea_column) VALUES ($1)", []byte{0x62, 0x79, 0x74, 0x65, 0x61})
		assert.NoError(t, err)
	})

	t.Run("Bytea型のデータが正常に取得できること", func(t *testing.T) {
		var row pq.ByteaArray
		err := tx.SelectContext(ctx, &row, "SELECT bytea_column FROM bytea_table")
		assert.NoError(t, err)

		assert.Equal(t, []byte("bytea"), row[0])
	})

	t.Run("Bytea型のデータが正常に更新できること", func(t *testing.T) {
		_, err := tx.ExecContext(ctx, "UPDATE bytea_table SET bytea_column = $1", []byte("bytea"))
		assert.NoError(t, err)
	})

	t.Run("Bytea型のデータが正常に削除できること", func(t *testing.T) {
		_, err := tx.ExecContext(ctx, "DELETE FROM bytea_table")
		assert.NoError(t, err)
	})
}
