package tasks

import (
	"context"
	"time"
	"log/slog"

	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/pkg/sharedtime"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/tasks/internal"

	"github.com/bwmarrin/discordgo"
)

func Run(ctx context.Context, dbv1 db.Driver, discord *discordgo.Session) {
	var messages []*discordgo.Message
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
					messages, err = internal.YoutubeRssReader(ctx, discord, repo, *webhook)
					if err != nil {
						slog.ErrorContext(ctx, "youtubeのwebhookの投稿に失敗しました。", "エラー", err.Error())
					}
				case "niconico":
					messages, err = internal.NiconicoRssReader(ctx, discord, repo, *webhook)
					if err != nil {
						slog.ErrorContext(ctx, "niconicoのwebhookの投稿に失敗しました。", "エラー", err.Error())
					}
				}
				if len(messages) > 0 {
					webhooks, err = repo.GetAllColumnsWebhooks(ctx)
					if err != nil {
						slog.ErrorContext(ctx, "webhookの取得に失敗しました。", "エラー", err.Error())
					}
				}
			}
			messages = nil
			for _, connect := range discord.VoiceConnections {
				connectTime := sharedtime.GetSharedTime(connect.GuildID)
				if time.Since(connectTime) > 10*time.Minute {
					err = connect.Disconnect()
					if err != nil {
						slog.ErrorContext(ctx, "ボイスチャンネルからの切断に失敗しました。", "エラー", err.Error())
					}
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
