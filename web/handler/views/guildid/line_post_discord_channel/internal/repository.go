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

func (r *Repository) GetLineChannel(ctx context.Context, channelID string) (LineChannel, error) {
	var channel LineChannel
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

func (r *Repository) InsertLineChannel(ctx context.Context, channelID string, guildID string) error {
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

func (r *Repository) GetLineNgType(ctx context.Context, channelID string) ([]int, error) {
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
