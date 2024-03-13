package internal

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

func (r *Repository) GetLineBots(ctx context.Context) ([]*LineBot, error) {
	var lineBots []*LineBot
	query := `
		SELECT
			guild_id,
			line_notify_token,
			line_bot_token,
			line_bot_secret,
			line_group_id,
			line_client_id,
			line_client_secret,
			default_channel_id,
			debug_mode
		FROM
			line_bot
		WHERE
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
			line_client_secret IS NOT NULL
	`
	err := r.db.SelectContext(ctx, &lineBots, query)
	return lineBots, err
}

func (r *Repository) GetLineBot(ctx context.Context, guildId string) (LineBot, error) {
	var lineBot LineBot
	query := `
		SELECT
			guild_id,
			line_notify_token,
			line_bot_token,
			line_bot_secret,
			line_group_id,
			line_client_id,
			line_client_secret,
			default_channel_id,
			debug_mode
		FROM
			line_bot
		WHERE
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
			line_client_secret IS NOT NULL
		AND
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineBot, query, guildId)
	return lineBot, err
}

func (r *Repository) GetLineBotIv(ctx context.Context, guildID string) (LineBotIv, error) {
	var lineBotIv LineBotIv
	query := `
		SELECT
			line_notify_token_iv,
			line_bot_token_iv,
			line_bot_secret_iv,
			line_client_id_iv,
			line_client_secret_iv,
			line_group_id_iv
		FROM
			line_bot_iv
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineBotIv, query, guildID)
	return lineBotIv, err
}
