package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type LineWorksBot struct {
	GuildID               string        `db:"guild_id"`
	LineWorksBotToken     pq.ByteaArray `db:"line_works_bot_token"`
	LineWorksRefreshToken pq.ByteaArray `db:"line_works_refresh_token"`
	LineWorksGroupID      pq.ByteaArray `db:"line_works_group_id"`
	LineWorksBotID        pq.ByteaArray `db:"line_works_bot_id"`
	LineWorksBotSecret    pq.ByteaArray `db:"line_works_bot_secret"`
	RefreshTokenExpiresAt pq.NullTime   `db:"refresh_token_expires_at"`
	DefaultChannelID      string        `db:"default_channel_id"`
	DebugMode             bool          `db:"debug_mode"`
}

type LineWorksBotNotClient struct {
	LineWorksBotToken     pq.ByteaArray `db:"line_works_bot_token"`
	LineWorksRefreshToken pq.ByteaArray `db:"line_works_refresh_token"`
	LineWorksGroupID      pq.ByteaArray `db:"line_works_group_id"`
	LineWorksBotID        pq.ByteaArray `db:"line_works_bot_id"`
	LineWorksBotSecret    pq.ByteaArray `db:"line_works_bot_secret"`
	RefreshTokenExpiresAt pq.NullTime   `db:"refresh_token_expires_at"`
	DefaultChannelID      string        `db:"default_channel_id"`
	DebugMode             bool          `db:"debug_mode"`
}

type LineWorksBotDefaultChannelID struct {
	GuildID          string `db:"guild_id"`
	DefaultChannelID string `db:"default_channel_id"`
}

func NewLineWorksBot(
	guildID string,
	lineWorksBotToken pq.ByteaArray,
	lineWorksRefreshToken pq.ByteaArray,
	lineWorksGroupID pq.ByteaArray,
	lineWorksBotID pq.ByteaArray,
	lineWorksBotSecret pq.ByteaArray,
	refreshTokenExpiresAt pq.NullTime,
	defaultChannelID string,
	debugMode bool,
) *LineWorksBot {
	return &LineWorksBot{
		GuildID:               guildID,
		LineWorksBotToken:     lineWorksBotToken,
		LineWorksRefreshToken: lineWorksRefreshToken,
		LineWorksGroupID:      lineWorksGroupID,
		LineWorksBotID:        lineWorksBotID,
		LineWorksBotSecret:    lineWorksBotSecret,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		DefaultChannelID:      defaultChannelID,
		DebugMode:             debugMode,
	}
}

func (r *Repository) GetAllLineWorksBots(ctx context.Context) ([]*LineWorksBot, error) {
	var lineWorksBots []*LineWorksBot
	query := `
		SELECT
			*
		FROM
			line_works_bot
	`
	err := r.db.SelectContext(ctx, &lineWorksBots, query)
	return lineWorksBots, err
}

func (r *Repository) InsertLineWorksBot(ctx context.Context, lineWorksBot *LineWorksBot) error {
	query := `
		INSERT INTO
			line_works_bot
			(
				guild_id,
				line_works_bot_token,
				line_works_refresh_token,
				line_works_group_id,
				line_works_bot_id,
				line_works_bot_secret,
				refresh_token_expires_at,
				default_channel_id,
				debug_mode
			)
		VALUES
			(
				:guild_id,
				:line_works_bot_token,
				:line_works_refresh_token,
				:line_works_group_id,
				:line_works_bot_id,
				:line_works_bot_secret,
				:refresh_token_expires_at,
				:default_channel_id,
				:debug_mode
			)
	`
	_, err := r.db.NamedExecContext(ctx, query, lineWorksBot)
	return err
}

func (r *Repository) GetLineWorksBotByGuildID(ctx context.Context, guildID string) (*LineWorksBot, error) {
	var lineWorksBot LineWorksBot
	query := `
		SELECT
			*
		FROM
			line_works_bot
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineWorksBot, query, guildID)
	return &lineWorksBot, err
}

func (r *Repository) GetLineWorksBotNotClientByGuildID(ctx context.Context, guildID string) (*LineWorksBotNotClient, error) {
	var lineWorksBot LineWorksBotNotClient
	query := `
		SELECT
			line_works_bot_token,
			line_works_refresh_token,
			line_works_group_id,
			line_works_bot_id,
			line_works_bot_secret,
			refresh_token_expires_at,
			default_channel_id,
			debug_mode
		FROM
			line_works_bot
		WHERE
			guild_id = $1
		AND
			line_works_bot_token IS NOT NULL
		AND
			line_works_refresh_token IS NOT NULL
		AND
			line_works_group_id IS NOT NULL
		AND
			line_works_bot_id IS NOT NULL
		AND
			line_works_bot_secret IS NOT NULL
	`
	err := r.db.GetContext(ctx, &lineWorksBot, query, guildID)
	return &lineWorksBot, err
}

func (r *Repository) GetLineWorksBotDefaultChannelIDByGuildID(ctx context.Context, guildID string) (*LineWorksBotDefaultChannelID, error) {
	var lineWorksBot LineWorksBotDefaultChannelID
	query := `
		SELECT
			default_channel_id
		FROM
			line_works_bot
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineWorksBot, query, guildID)
	return &lineWorksBot, err
}

func (r *Repository) UpdateLineWorksBot(ctx context.Context, lineWorksBot *LineWorksBot) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineWorksBot.LineWorksBotToken) > 0 {
		setQueryArray = append(setQueryArray, "line_works_bot_token = :line_works_bot_token")
	}
	if len(lineWorksBot.LineWorksRefreshToken) > 0 {
		setQueryArray = append(setQueryArray, "line_works_refresh_token = :line_works_refresh_token")
	}
	if len(lineWorksBot.LineWorksGroupID) > 0 {
		setQueryArray = append(setQueryArray, "line_works_group_id = :line_works_group_id")
	}
	if len(lineWorksBot.LineWorksBotID) > 0 {
		setQueryArray = append(setQueryArray, "line_works_bot_id = :line_works_bot_id")
	}
	if len(lineWorksBot.LineWorksBotSecret) > 0 {
		setQueryArray = append(setQueryArray, "line_works_bot_secret = :line_works_bot_secret")
	}
	if lineWorksBot.RefreshTokenExpiresAt.Valid {
		setQueryArray = append(setQueryArray, "refresh_token_expires_at = :refresh_token_expires_at")
	}
	if lineWorksBot.DefaultChannelID != "" {
		setQueryArray = append(setQueryArray, "default_channel_id = :default_channel_id")
	}
	setQueryArray = append(setQueryArray, "debug_mode = :debug_mode")
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
		return nil
	}
	query := fmt.Sprintf(`
		UPDATE
			line_works_bot
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)
	_, err := r.db.NamedExecContext(ctx, query, lineWorksBot)
	return err
}
