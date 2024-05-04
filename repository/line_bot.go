package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type LineBot struct {
	GuildID          string        `db:"guild_id"`
	LineNotifyToken  pq.ByteaArray `db:"line_notify_token"`
	LineBotToken     pq.ByteaArray `db:"line_bot_token"`
	LineBotSecret    pq.ByteaArray `db:"line_bot_secret"`
	LineGroupID      pq.ByteaArray `db:"line_group_id"`
	LineClientID     pq.ByteaArray `db:"line_client_id"`
	LineClientSecret pq.ByteaArray `db:"line_client_secret"`
	DefaultChannelID string        `db:"default_channel_id"`
	DebugMode        bool          `db:"debug_mode"`
}

type LineBotNotClient struct {
	LineNotifyToken  pq.ByteaArray `db:"line_notify_token"`
	LineBotToken     pq.ByteaArray `db:"line_bot_token"`
	LineBotSecret    pq.ByteaArray `db:"line_bot_secret"`
	LineGroupID      pq.ByteaArray `db:"line_group_id"`
	DefaultChannelID string        `db:"default_channel_id"`
	DebugMode        bool          `db:"debug_mode"`
}

type LineBotDefaultChannelID struct {
	GuildID          string `db:"guild_id"`
	DefaultChannelID string `db:"default_channel_id"`
}

func NewLineBot(
	guildID string,
	lineNotifyToken pq.ByteaArray,
	lineBotToken pq.ByteaArray,
	lineBotSecret pq.ByteaArray,
	lineGroupID pq.ByteaArray,
	lineClientID pq.ByteaArray,
	lineClientSecret pq.ByteaArray,
	defaultChannelID string,
	debugMode bool,
) *LineBot {
	return &LineBot{
		GuildID:          guildID,
		LineNotifyToken:  lineNotifyToken,
		LineBotToken:     lineBotToken,
		LineBotSecret:    lineBotSecret,
		LineGroupID:      lineGroupID,
		LineClientID:     lineClientID,
		LineClientSecret: lineClientSecret,
		DefaultChannelID: defaultChannelID,
		DebugMode:        debugMode,
	}
}

func (r *Repository) InsertLineBot(ctx context.Context, lineBot *LineBot) error {
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

func (r *Repository) GetAllColumnsLineBots(ctx context.Context) ([]*LineBot, error) {
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

func (r *Repository) GetAllColumnsLineBot(ctx context.Context, guildId string) (LineBot, error) {
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
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineBot, query, guildId)
	return lineBot, err
}

func (r *Repository) GetLineBotDefaultChannelID(ctx context.Context, guildID string) (LineBotDefaultChannelID, error) {
	var lineBot LineBotDefaultChannelID
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

func (r *Repository) GetLineBotNotClient(ctx context.Context, guildID string) (LineBotNotClient, error) {
	var lineBot LineBotNotClient
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

func (r *Repository) UpdateLineBot(ctx context.Context, lineBot *LineBot) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineBot.LineNotifyToken) > 0 && len(lineBot.LineNotifyToken[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_notify_token = :line_notify_token")
	}
	if len(lineBot.LineBotToken) > 0 && len(lineBot.LineBotToken[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_token = :line_bot_token")
	}
	if len(lineBot.LineBotSecret) > 0 && len(lineBot.LineBotSecret[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_secret = :line_bot_secret")
	}
	if len(lineBot.LineGroupID) > 0 && len(lineBot.LineGroupID[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_group_id = :line_group_id")
	}
	if len(lineBot.LineClientID) > 0 && len(lineBot.LineClientID[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_client_id = :line_client_id")
	}
	if len(lineBot.LineClientSecret) > 0 && len(lineBot.LineClientSecret[0]) > 0 {
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
		fmt.Println("No update value")
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
