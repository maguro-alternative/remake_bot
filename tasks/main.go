package tasks

import (
	"context"
	"time"

	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/bwmarrin/discordgo"
)

func Run(ctx context.Context, dbv1 db.Driver, discord *discordgo.Session) error {
	// ここにタスクを書く
	oneMinute := time.NewTicker(1 * time.Minute)
	tenMinute := time.NewTicker(10 * time.Minute)
	repo := repository.NewRepository(dbv1)
	webhooks, err := repo.GetAllColumnsWebhooks(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case <-oneMinute.C:
			for _, webhook := range webhooks {
				switch webhook.SubscriptionType {
				case "youtube":
					// ここにタスクを書く
				case "niconico":
					// ここにタスクを書く
				}
			}
		case <-tenMinute.C:
			webhooks, err = repo.GetAllColumnsWebhooks(ctx)
			if err != nil {
				return err
			}
		}
	}
}
