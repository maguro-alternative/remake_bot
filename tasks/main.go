package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

func Run(ctx context.Context, dbv1 db.Driver,) error {
	// ここにタスクを書く
	oneMinute := time.NewTicker(1 * time.Minute)
	tenMinute := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-oneMinute.C:
			fmt.Println("1分経過")
		case <-tenMinute.C:
			fmt.Println("10分経過")
		}
	}
	return nil
}
