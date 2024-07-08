// sharedパッケージ
package sharedtime

import (
    "sync"
    "time"
)

var (
    sharedTimes map[string]time.Time
    lock       sync.Mutex
)

func SetSharedTime(t map[string]time.Time) {
    lock.Lock()
    defer lock.Unlock()
    sharedTimes = t
}

func GetSharedTime(guildId string) time.Time {
    lock.Lock()
    defer lock.Unlock()
    return sharedTimes[guildId]
}
