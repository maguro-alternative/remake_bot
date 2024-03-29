package on_message_create

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type Repository struct {
	db db.Driver
}

func NewRepository(db db.Driver) *Repository {
	return &Repository{
		db: db,
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

func (r *Repository) GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error) {
	var ngTypes []int
	query := `
		SELECT
			type
		FROM
			line_ng_discord_message_type
		WHERE
			channel_id = $1
	`
	err := r.db.SelectContext(ctx, &ngTypes, query, channelID)
	return ngTypes, err
}

func (r *Repository) GetLineNgDiscordID(ctx context.Context, channelID string) ([]LineNgID, error) {
	var ngIDs []LineNgID
	query := `
		SELECT
			id,
			id_type
		FROM
			line_ng_discord_id
		WHERE
			channel_id = $1
	`
	err := r.db.SelectContext(ctx, &ngIDs, query, channelID)
	return ngIDs, err
}

func (r *Repository) GetLineBot(ctx context.Context, guildID string) (LineBot, error) {
	var lineBot LineBot
	query := `
		SELECT
			line_notify_token,
			line_bot_token,
			line_bot_secret,
			line_group_id,
			default_channel_id,
			debug_mode
		FROM
			line_bot
		WHERE
			guild_id = $1
		AND
			line_notify_token IS NOT NULL
		AND
			line_bot_token IS NOT NULL
		AND
			line_bot_secret IS NOT NULL
		AND
			line_group_id IS NOT NULL
	`
	err := r.db.GetContext(ctx, &lineBot, query, guildID)
	return lineBot, err
}

func (r *Repository) GetLineBotIv(ctx context.Context, guildID string) (LineBotIv, error) {
	var lineBotIv LineBotIv
	query := `
		SELECT
			line_notify_token_iv,
			line_bot_token_iv,
			line_bot_secret_iv,
			line_group_id_iv
		FROM
			line_bot_iv
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineBotIv, query, guildID)
	return lineBotIv, err
}
