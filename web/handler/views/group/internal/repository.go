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

func (r *Repository) GetLineBot(ctx context.Context, guildID string) (LineBot, error) {
	var lineBot LineBot
	query := `
		SELECT
			default_channel_id
		FROM
			line_bot
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineBot, query, guildID)
	return lineBot, err
}
