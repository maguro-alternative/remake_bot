package internal

import (
	"context"
	"fmt"
	"reflect"
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
	structTypeOf := reflect.TypeOf(lineBotIv)

	// 受け取った構造体のフィールドのみを更新する
	for i := 0; i < structTypeOf.NumField(); i++ {
		field := structTypeOf.Field(i).Tag.Get("db")
		if field == "" || field == "guild_id"{
			continue
		}
		if i == structTypeOf.NumField()-1 {
			setNameQuery += field + " = :" + field
			continue
		}
		setNameQuery += field + " = :" + field + ","
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
