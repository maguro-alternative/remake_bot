package repository

import (
	"context"
	"fmt"
	"strings"
)

type LineNgDiscordMessageType struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	Type      int    `db:"type"`
}

func (r *Repository) InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgTypes []LineNgDiscordMessageType) error {
	query := `
		INSERT INTO line_ng_discord_message_type (
			channel_id,
			guild_id,
			type
		) VALUES (
			:channel_id,
			:guild_id,
			:type
		) ON CONFLICT (channel_id, type) DO NOTHING
	`
	for _, lineNgType := range lineNgTypes {
		_, err := r.db.NamedExecContext(ctx, query, lineNgType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgTypes []LineNgDiscordMessageType) error {
	var values []string
	for _, lineNgType := range lineNgTypes {
		values = append(values, fmt.Sprintf("('%s', '%s', %d)", lineNgType.ChannelID, lineNgType.GuildID, lineNgType.Type))
		_, err := r.db.ExecContext(ctx, "DELETE FROM line_ng_discord_message_type WHERE channel_id = $1", lineNgType.ChannelID)
		if err != nil {
			return err
		}
	}
	if len(values) == 0 {
		return nil
	}
	// INSERT されるもの以外を削除
	query := fmt.Sprintf(`
		INSERT INTO line_ng_discord_message_type (
			channel_id,
			guild_id,
			type
		) VALUES
				%s
		ON CONFLICT (channel_id, type) DO NOTHING
	`, strings.Join(values, ","))
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *Repository) GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error) {
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


