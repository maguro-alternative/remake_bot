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

	if len(lineBot.LineNotifyToken) > 0 {
		setQueryArray = append(setQueryArray, "line_notify_token = :line_notify_token")
	}
	if len(lineBot.LineBotToken) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_token = :line_bot_token")
	}
	if len(lineBot.LineBotSecret) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_secret = :line_bot_secret")
	}
	if len(lineBot.LineGroupID) > 0 {
		setQueryArray = append(setQueryArray, "line_group_id = :line_group_id")
	}
	if len(lineBot.LineClientID) > 0 {
		setQueryArray = append(setQueryArray, "line_client_id = :line_client_id")
	}
	if len(lineBot.LineClientSecret) > 0 {
		setQueryArray = append(setQueryArray, "line_client_secret = :line_client_secret")
	}
	if lineBot.DefaultChannelID != "" {
		setQueryArray = append(setQueryArray, "default_channel_id = :default_channel_id")
	}
	if lineBot.DebugMode {
		setQueryArray = append(setQueryArray, "debug_mode = :debug_mode")
	}
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
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

func (r *Repository) UpdateLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineBotIv.LineNotifyTokenIv) > 0 {
		setQueryArray = append(setQueryArray, "line_notify_token_iv = :line_notify_token_iv")
	}
	if len(lineBotIv.LineBotTokenIv) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_token_iv = :line_bot_token_iv")
	}
	if len(lineBotIv.LineBotSecretIv) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_secret_iv = :line_bot_secret_iv")
	}
	if len(lineBotIv.LineGroupIDIv) > 0 {
		setQueryArray = append(setQueryArray, "line_group_id_iv = :line_group_id_iv")
	}
	if len(lineBotIv.LineClientIDIv) > 0 {
		setQueryArray = append(setQueryArray, "line_client_id_iv = :line_client_id_iv")
	}
	if len(lineBotIv.LineClientSecretIv) > 0 {
		setQueryArray = append(setQueryArray, "line_client_secret_iv = :line_client_secret_iv")
	}
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE
			line_bot_iv
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)
	_, err := r.db.NamedExecContext(ctx, query, lineBotIv)
	return err
}

func (r *Repository) InsertLineBot(ctx context.Context, lineBot *LineBot)  error {
	query := `
		INSERT INTO line_bot (
			guild_id,
			line_notify_token,
			line_bot_token,
			line_bot_secret,
			line_group_id,
			line_client_id,
			line_client_secret,
			default_channel_id,
			debug_mode
		) VALUES (
			:guild_id,
			:line_notify_token,
			:line_bot_token,
			:line_bot_secret,
			:line_group_id,
			:line_client_id,
			:line_client_secret,
			:default_channel_id,
			:debug_mode
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, lineBot)
	return err
}

func (r *Repository) InsertLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error {
	query := `
		INSERT INTO line_bot_iv (
			guild_id,
			line_notify_token_iv,
			line_bot_token_iv,
			line_bot_secret_iv,
			line_group_id_iv,
			line_client_id_iv,
			line_client_secret_iv
		) VALUES (
			:guild_id,
			:line_notify_token_iv,
			:line_bot_token_iv,
			:line_bot_secret_iv,
			:line_group_id_iv,
			:line_client_id_iv,
			:line_client_secret_iv
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, lineBotIv)
	return err
}
