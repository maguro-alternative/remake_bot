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
			line_notify_token,
			line_bot_token,
			line_bot_secret,
			line_group_id,
			line_client_id,
			line_client_sercret,
			iv,
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
			line_client_sercret IS NOT NULL
		AND
			iv IS NOT NULL
	`
	err := r.db.SelectContext(ctx, &lineBots, query)
	return lineBots, err
}
