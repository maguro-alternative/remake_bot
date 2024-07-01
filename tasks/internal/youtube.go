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

func Run(
	ctx context.Context,
	repo repository.RepositoryFunc,
	webhook repository.Webhook,
	discordSession mock.Session,
) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message
	var mentionsMessage string
	var lastPostedAt time.Time
	// ここにタスクを書く
	feed, err := gofeed.NewParser().ParseURL(fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", webhook.SubscriptionID))
    if err != nil {
		return nil, err
	}
	for _, item := range feed.Items {
        if item == nil {
            break
        }
		if webhook.LastPostedAt.After(*item.PublishedParsed) {
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
			message, err := discordSession.WebhookExecute(w.ID, w.Token, false, &discordgo.WebhookParams{
				Content: fmt.Sprintf("%s\n%s\n%s", mentionsMessage, item.Title, item.Link),
			})
			if err != nil {
				return nil, err
			}
			if lastPostedAt.Before(*item.PublishedParsed) {
				lastPostedAt = *item.PublishedParsed
			}
			messages = append(messages, message)
		}
	}
	if !lastPostedAt.IsZero() {
		err = repo.UpdateWebhookWithLastPostedAt(ctx, *webhook.WebhookSerialID, lastPostedAt)
		if err != nil {
			return nil, err
		}
	}
	return messages, nil
}
