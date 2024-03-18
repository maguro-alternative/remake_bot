package internal

import (
	"context"
	"fmt"
	"strings"

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

func (r *Repository) UpdateLineBot(ctx context.Context, lineBot *LineBot) error {
	var setNameQuery string
	var setQueryArray []string

	if lineBot.DefaultChannelID != "" {
		setQueryArray = append(setQueryArray, "default_channel_id = :default_channel_id")
	}
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
		fmt.Println("setNameQuery is empty")
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE
			line_bot
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)
	_, err := r.db.NamedExecContext(ctx, query, lineBot)
	return err
}
