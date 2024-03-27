package repository

import (
	"context"
)

type LinePostDiscordChannel struct {
	Ng         bool `db:"ng"`
	BotMessage bool `db:"bot_message"`
}

type LinePostDiscordChannelAllColumns struct {
	ChannelID  string `db:"channel_id"`
	GuildID    string `db:"guild_id"`
	Ng         bool   `db:"ng"`
	BotMessage bool   `db:"bot_message"`
}

func NewLinePostDiscordChannel(channelID, guildID string, ng, botMessage bool) *LinePostDiscordChannelAllColumns {
	return &LinePostDiscordChannelAllColumns{
		ChannelID:  channelID,
		GuildID:    guildID,
		Ng:         ng,
		BotMessage: botMessage,
	}
}

func (r *Repository) GetLinePostDiscordChannel(ctx context.Context, channelID string) (LinePostDiscordChannel, error) {
	var channel LinePostDiscordChannel
	query := `
		SELECT
			ng,
			bot_message
		FROM
			line_post_discord_channel
		WHERE
			channel_id = $1
	`
	err := r.db.GetContext(ctx, &channel, query, channelID)
	return channel, err
}

func (r *Repository) UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error {
	query := `
		UPDATE
			line_post_discord_channel
		SET
			ng = :ng,
			bot_message = :bot_message
		WHERE
			channel_id = :channel_id
	`
	_, err := r.db.NamedExecContext(ctx, query, lineChannel)
	return err
}

func (r *Repository) InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error {
	query := `
		INSERT INTO line_post_discord_channel (
			channel_id,
			guild_id,
			ng,
			bot_message
		) VALUES (
			$1,
			$2,
			false,
			false
		)
	`
	_, err := r.db.ExecContext(ctx, query, channelID, guildID)
	return err
}