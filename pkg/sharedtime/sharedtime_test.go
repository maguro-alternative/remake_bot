package sharedtime

import (
    "reflect"
    "testing"
    "time"
)

func TestSetAndGetSharedTime(t *testing.T) {
    // テスト用の時間データを準備
    testTime := map[string]time.Time{
        "guild1": time.Now(),
    }

    // SetSharedTimeを呼び出して、sharedTimesに値を設定
    SetSharedTime(testTime)

    // GetSharedTimeを呼び出して、設定した値を取得
    gotTime := GetSharedTime("guild1")

    // 設定した時間と取得した時間が等しいか検証
    if !reflect.DeepEqual(testTime["guild1"], gotTime) {
        t.Errorf("GetSharedTime() = %v, want %v", gotTime, testTime["guild1"])
    }
}

func TestGetSharedTimeNonExistentKey(t *testing.T) {
    // 存在しないキーでGetSharedTimeを呼び出す
    gotTime := GetSharedTime("nonexistent")

    // ゼロ値のtime.Timeが返されることを検証
    if !gotTime.IsZero() {
        t.Errorf("GetSharedTime() with nonexistent key should return zero time, got %v", gotTime)
    }
}
