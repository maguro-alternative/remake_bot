package internal

import (
	"context"
	"time"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/stretchr/testify/assert"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

func TestYoutubeRssReader(t *testing.T) {
	ctx := context.Background()
	discordSession := &mock.SessionMock{
		WebhookFunc: func(webhookID string, options ...discordgo.RequestOption) (*discordgo.Webhook, error) {
			return &discordgo.Webhook{}, nil
		},
		WebhookExecuteFunc: func(webhookID string, token string, wait bool, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (*discordgo.Message, error) {
			return &discordgo.Message{}, nil
		},
	}
	repo := &repository.Repository{}
	webhook := repository.Webhook{}
	feed := &gofeed.Feed{
		Items: []*gofeed.Item{
			{
				PublishedParsed: &time.Time{},
			},
		},
	}
	t.Run("YoutubeのRss取得に成功すること", func(t *testing.T) {
		_, err := run(ctx, discordSession, repo, webhook, feed)
		assert.NoError(t, err)
	})
}
