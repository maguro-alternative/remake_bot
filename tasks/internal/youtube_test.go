package internal

import (
	"context"
	"testing"
	"time"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/stretchr/testify/assert"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

func TestYoutubeRssReader(t *testing.T) {
	ctx := context.Background()
	beforePostAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	afterPostAt := time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
	discordSession := &mock.SessionMock{
		WebhookFunc: func(webhookID string, options ...discordgo.RequestOption) (*discordgo.Webhook, error) {
			return &discordgo.Webhook{}, nil
		},
		WebhookExecuteFunc: func(webhookID string, token string, wait bool, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (*discordgo.Message, error) {
			return &discordgo.Message{}, nil
		},
	}
	repo := &repository.RepositoryFuncMock{
		GetWebhookUserMentionWithWebhookSerialIDFunc: func(ctx context.Context, webhookSerialID int64) ([]*repository.WebhookUserMention, error) {
			return []*repository.WebhookUserMention{}, nil
		},
		GetWebhookRoleMentionWithWebhookSerialIDFunc: func(ctx context.Context, webhookSerialID int64) ([]*repository.WebhookRoleMention, error) {
			return []*repository.WebhookRoleMention{}, nil
		},
	}
	webhookSerialId := int64(1)
	webhook := repository.Webhook{
		WebhookSerialID:  &webhookSerialId,
		GuildID:          "1111",
		WebhookID:        "2222",
		SubscriptionType: "youtube",
		SubscriptionID:   "test",
		LastPostedAt:     beforePostAt,
	}
	feed := &gofeed.Feed{
		Items: []*gofeed.Item{
			{
				Title:           "test",
				Link:            "https://www.youtube.com/watch?v=test",
				PublishedParsed: &afterPostAt,
			},
		},
	}
	t.Run("YoutubeのRss取得に成功すること", func(t *testing.T) {
		messages, err := run(ctx, discordSession, repo, webhook, feed)
		assert.NoError(t, err)
		assert.Len(t, messages, 1)
	})
}
