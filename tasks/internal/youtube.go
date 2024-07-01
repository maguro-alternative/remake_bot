package internal

import (
	"context"
	"fmt"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

func Run(
	ctx context.Context,
	webhook repository.Webhook,
	discordSession *discordgo.Session,
) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message
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
			message, err := discordSession.WebhookExecute(w.ID, w.Token, false, &discordgo.WebhookParams{
				Content: fmt.Sprintf("%s\n%s", item.Title, item.Link),
			})
			if err != nil {
				return nil, err
			}
			messages = append(messages, message)
		}
	}
	return messages, nil
}
