package internal

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/mmcdole/gofeed"
)

func Run(
	ctx context.Context,
	youtubeId string,
	discord *discordgo.Session,
) error {
	// ここにタスクを書く
	feed, err := gofeed.NewParser().ParseURL(fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", youtubeId))
    if err != nil {
		return err
	}
	for _, item := range feed.Items {
        if item == nil {
            break
        }
	}
	return nil
}
