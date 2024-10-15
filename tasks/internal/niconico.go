package internal

import (
	"context"
	"fmt"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

func NiconicoRssReader(ctx context.Context,
	discordSession mock.Session,
	repo repository.RepositoryFunc,
	webhook repository.Webhook,
) ([]*discordgo.Message, error) {
	feed, err := gofeed.NewParser().ParseURL(fmt.Sprintf("https://www.nicovideo.jp/user/%s/video?rss=2.0", webhook.SubscriptionID))
	if err != nil {
		return nil, err
	}
	return run(ctx, discordSession, repo, webhook, feed)
}

