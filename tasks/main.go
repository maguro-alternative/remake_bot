package tasks

import (
	"context"
	"time"
	"log/slog"

	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/tasks/internal"

	"github.com/bwmarrin/discordgo"
)

func Run(ctx context.Context, dbv1 db.Driver, discord *discordgo.Session) {
	// ここにタスクを書く
	oneMinute := time.NewTicker(1 * time.Minute)
	tenMinute := time.NewTicker(10 * time.Minute)
	repo := repository.NewRepository(dbv1)
	webhooks, err := repo.GetAllColumnsWebhooks(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "webhookの取得に失敗しました。", "エラー", err.Error())
	}
	for {
		select {
		case <-oneMinute.C:
			for _, webhook := range webhooks {
				switch webhook.SubscriptionType {
				case "youtube":
					_, err := internal.YoutubeRssReader(ctx, discord, repo, *webhook)
					if err != nil {
						slog.ErrorContext(ctx, "youtubeのwebhookの投稿に失敗しました。", "エラー", err.Error())
					}
				case "niconico":
					// Todo: ニコニコ動画のRSSリーダーを実装する
				}
			}
		case <-tenMinute.C:
			webhooks, err = repo.GetAllColumnsWebhooks(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "webhookの取得に失敗しました。", "エラー", err.Error())
			}
		}
	}
}
