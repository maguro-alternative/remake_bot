// sharedパッケージ
package sharedtime

import (
	"sync"
	"time"
)

var (
	sharedTimes map[string]time.Time = make(map[string]time.Time)
	lock        sync.Mutex
)

func SetSharedTime(guildId string, t time.Time) {
	lock.Lock()
	defer lock.Unlock()
	sharedTimes[guildId] = t
}

func GetSharedTime(guildId string) time.Time {
	lock.Lock()
	defer lock.Unlock()
	return sharedTimes[guildId]
}
