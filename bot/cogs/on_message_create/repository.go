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

func (r *Repository) GetLineChannel(ctx context.Context, channelID string) (LineChannel, error) {
	var channel LineChannel
	query := `
		SELECT
			ng,
			bot_message
		FROM
			line_channels
		WHERE
			channel_id = $1
	`
	err := r.db.GetContext(ctx, &channel, query, channelID)
	return channel, err
}

func (r *Repository) InsertLineChannel(ctx context.Context, channelID string, guildID string) error {
	query := `
		INSERT INTO line_channels (
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

func (r *Repository) GetLineNgType(ctx context.Context, guildID string) ([]string, error) {
	var ngTypes []string
	query := `
		SELECT
			type
		FROM
			line_ng_types
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &ngTypes, query, guildID)
	return ngTypes, err
}

func (r *Repository) GetLineBot(ctx context.Context, guildID string) (LineBot, error) {
	var lineBot LineBot
	query := `
		SELECT
			line_notify_token,
			line_bot_token,
			line_bot_secret,
			line_group_id,
			iv,
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
		AND
			line_client_id IS NOT NULL
		AND
			line_client_sercret IS NOT NULL
		AND
			iv IS NOT NULL
	`
	err := r.db.GetContext(ctx, &lineBot, query, guildID)
	return lineBot, err
}
