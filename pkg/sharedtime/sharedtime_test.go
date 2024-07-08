package sharedtime

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestSetAndGetSharedTime(t *testing.T) {
    t.Run("SetSharedTimeで設定した値をGetSharedTimeで取得できる", func(t *testing.T) {
        // SetSharedTimeで設定する値
        wantTime := time.Now()

        // SetSharedTimeで値を設定
        SetSharedTime("key", wantTime)

        // GetSharedTimeで値を取得
        gotTime := GetSharedTime("key")

        // SetSharedTimeで設定した値とGetSharedTimeで取得した値が一致することを検証
        assert.Equal(t, wantTime, gotTime)
    })
}

func TestGetSharedTimeNonExistentKey(t *testing.T) {
    t.Run("存在しないキーでGetSharedTimeを実行するとゼロ値が返る", func(t *testing.T) {
        // GetSharedTimeで存在しないキーを指定
        gotTime := GetSharedTime("non-existent-key")

        // ゼロ値が返ることを検証
        assert.Equal(t, time.Time{}, gotTime)
    })
}
