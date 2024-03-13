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

func (r *Repository) InsertLineBot(ctx context.Context, lineBot *LineBot)  error {
	query := `
		INSERT INTO line_bot (
			guild_id,
			default_channel_id,
			debug_mode
		) VALUES (
			:guild_id,
			:default_channel_id,
			:debug_mode
		) ON CONFLICT (guild_id) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, lineBot)
	return err
}

func (r *Repository) InsertLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error {
	query := `
		INSERT INTO line_bot_iv (
			guild_id
		) VALUES (
			:guild_id
		) ON CONFLICT (guild_id) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, lineBotIv)
	return err
}

func (r *Repository) GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error) {
	var code int64
	if permissionType == "" {
		return code, nil
	}
	query := `
		SELECT
			code
		FROM
			permissions_code
		WHERE
			guild_id = $1 AND
			type = $2
	`
	err := r.db.GetContext(ctx, &code, query, guildID, permissionType)
	return code, err
}

func (r *Repository) GetPermissionIDs(ctx context.Context, guildID, permissionType string) ([]PermissionID, error) {
	var permissionIDs []PermissionID
	query := `
		SELECT
			target_type,
			target_id,
			permission
		FROM
			permissions_id
		WHERE
			guild_id = $1 AND
			type = $2
	`
	err := r.db.SelectContext(ctx, &permissionIDs, query, guildID, permissionType)
	return permissionIDs, err
}
