package repository

import (
	"context"
	"fmt"
	"strings"
)

type LineNgDiscordUserIDAllCoulmns struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	UserID    string `db:"user_id"`
}

func NewLineNgDiscordUserID(channelID, guildID, userID string) *LineNgDiscordUserIDAllCoulmns {
	return &LineNgDiscordUserIDAllCoulmns{
		ChannelID: channelID,
		GuildID:   guildID,
		UserID:    userID,
	}
}

func (r *Repository) GetLineNgDiscordUserID(ctx context.Context, channelID string) ([]string, error) {
	var ngIDs []string
	query := `
		SELECT
			user_id
		FROM
			line_ng_discord_user_id
		WHERE
			channel_id = $1
	`
	err := r.db.SelectContext(ctx, &ngIDs, query, channelID)
	return ngIDs, err
}

func (r *Repository) InsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error {
	query := `
		INSERT INTO line_ng_discord_user_id (
			channel_id,
			guild_id,
			user_id
		) VALUES (
			:channel_id,
			:guild_id,
			:user_id
		) ON CONFLICT (channel_id, user_id) DO NOTHING
	`
	for _, lineNgID := range lineNgDiscordUserIDs {
		_, err := r.db.NamedExecContext(ctx, query, lineNgID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeleteNotInsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error {
	var values []string
	for _, lineNgType := range lineNgDiscordUserIDs {
		values = append(values, fmt.Sprintf("('%s', '%s', '%s')", lineNgType.ChannelID, lineNgType.GuildID, lineNgType.UserID))
		_, err := r.db.ExecContext(ctx, "DELETE FROM line_ng_discord_user_id WHERE channel_id = $1", lineNgType.ChannelID)
		if err != nil {
			return err
		}
	}
	if len(values) == 0 {
		return nil
	}
	// INSERT されるもの以外を削除
	query := fmt.Sprintf(`
		INSERT INTO line_ng_discord_user_id (
			channel_id,
			guild_id,
			user_id
		) VALUES
			%s
		ON CONFLICT (channel_id, user_id) DO NOTHING
	`, strings.Join(values, ","))
	_, err := r.db.ExecContext(ctx, query)
	return err
}
