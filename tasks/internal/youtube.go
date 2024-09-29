package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

func YoutubeRssReader(
	ctx context.Context,
	discordSession mock.Session,
	repo repository.RepositoryFunc,
	webhook repository.Webhook,
) ([]*discordgo.Message, error) {
	feed, err := gofeed.NewParser().ParseURL(fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", webhook.SubscriptionID))
	if err != nil {
		return nil, err
	}
	return run(ctx, discordSession, repo, webhook, feed)
}

func run(
	ctx context.Context,
	discordSession mock.Session,
	repo repository.RepositoryFunc,
	webhook repository.Webhook,
	feed *gofeed.Feed,
) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message
	var mentionsMessage string
	var lastPostedAt time.Time
	for _, item := range feed.Items {
		if item == nil {
			break
		}
		if webhook.LastPostedAt.Before(*item.PublishedParsed) {
			w, err := discordSession.Webhook(webhook.WebhookID)
			if err != nil {
				return nil, err
			}
			userMentions, err := repo.GetWebhookUserMentionWithWebhookSerialID(ctx, *webhook.WebhookSerialID)
			if err != nil {
				return nil, err
			}
			roleMentions, err := repo.GetWebhookRoleMentionWithWebhookSerialID(ctx, *webhook.WebhookSerialID)
			if err != nil {
				return nil, err
			}
			for _, userMention := range userMentions {
				mentionsMessage += fmt.Sprintf("<@%s> ", userMention.UserID)
			}
			for _, roleMention := range roleMentions {
				mentionsMessage += fmt.Sprintf("<@&%s> ", roleMention.RoleID)
			}
			if mentionsMessage != "" {
				mentionsMessage += "\n"
			}
			threads, err := repo.GetWebhookThreadWithWebhookSerialID(ctx, *webhook.WebhookSerialID)
			if err != nil {
				return nil, err
			}
			for _, thread := range threads {
				message, err := discordSession.WebhookThreadExecute(w.ID, w.Token, false, thread.ThreadID, &discordgo.WebhookParams{
					Content: fmt.Sprintf("%s%s\n%s", mentionsMessage, item.Title, item.Link),
				})
				if err != nil {
					return nil, err
				}
				messages = append(messages, message)
			}
			if len(threads) == 0 {
				fmt.Println("%#v",webhook)
				message, err := discordSession.WebhookExecute(w.ID, w.Token, false, &discordgo.WebhookParams{
					Content: fmt.Sprintf("%s%s\n%s", mentionsMessage, item.Title, item.Link),
				})
				if err != nil {
					return nil, err
				}
				messages = append(messages, message)
			}
			if lastPostedAt.Before(*item.PublishedParsed) {
				lastPostedAt = *item.PublishedParsed
			}
			userMentions = nil
			roleMentions = nil
			mentionsMessage = ""
		}
	}
	if !lastPostedAt.IsZero() {
		err := repo.UpdateWebhookWithLastPostedAt(ctx, *webhook.WebhookSerialID, lastPostedAt)
		if err != nil {
			return nil, err
		}
	}
	return messages, nil
}
