package repository

import (
	"context"
	"fmt"
	"strings"
)

type LineNgDiscordRoleIDAllCoulmns struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	ID        string `db:"id"`
}

func NewLineNgDiscordRoleID(channelID, guildID, id string) *LineNgDiscordRoleIDAllCoulmns {
	return &LineNgDiscordRoleIDAllCoulmns{
		ChannelID: channelID,
		GuildID:   guildID,
		ID:        id,
	}
}

func (r *Repository) GetLineNgDiscordRoleID(ctx context.Context, channelID string) ([]string, error) {
	var ngIDs []string
	query := `
		SELECT
			id
		FROM
			line_ng_discord_role_id
		WHERE
			channel_id = $1
	`
	err := r.db.GetContext(ctx, &ngIDs, query, channelID)
	return ngIDs, err
}

func (r *Repository) InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	query := `
		INSERT INTO line_ng_discord_role_id (
			channel_id,
			guild_id,
			id
		) VALUES (
			:channel_id,
			:guild_id,
			:id
		) ON CONFLICT (channel_id, id) DO NOTHING
	`
	for _, lineNgID := range lineNgDiscordRoleIDs {
		_, err := r.db.NamedExecContext(ctx, query, lineNgID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeleteNotInsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	var values []string
	for _, lineNgType := range lineNgDiscordRoleIDs {
		values = append(values, fmt.Sprintf("('%s', '%s', '%s')", lineNgType.ChannelID, lineNgType.GuildID, lineNgType.ID))
		_, err := r.db.ExecContext(ctx, "DELETE FROM line_ng_discord_role_id WHERE channel_id = $1", lineNgType.ChannelID)
		if err != nil {
			return err
		}
	}
	if len(values) == 0 {
		return nil
	}
	// INSERT されるもの以外を削除
	query := fmt.Sprintf(`
		INSERT INTO line_ng_discord_role_id (
			channel_id,
			guild_id,
			id
		) VALUES
			%s
		ON CONFLICT (channel_id, id) DO NOTHING
	`, strings.Join(values, ","))
	_, err := r.db.ExecContext(ctx, query)
	return err
}

