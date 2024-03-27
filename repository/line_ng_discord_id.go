package repository

import (
	"context"
	"fmt"
	"strings"
)

type LineNgDiscordID struct {
	ID     string `db:"id"`
	IDType string `db:"id_type"`
}

type LineNgDiscordIDAllCoulmns struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	ID        string `db:"id"`
	IDType    string `db:"id_type"`
}

func (r *Repository) GetLineNgDiscordID(ctx context.Context, channelID string) ([]LineNgDiscordID, error) {
	var ngIDs []LineNgDiscordID
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

func (r *Repository) InsertLineNgDiscordIDs(ctx context.Context, lineNgIDs []LineNgDiscordIDAllCoulmns) error {
	query := `
		INSERT INTO line_ng_discord_id (
			channel_id,
			guild_id,
			id,
			id_type
		) VALUES (
			:channel_id,
			:guild_id,
			:id,
			:id_type
		) ON CONFLICT (channel_id, id) DO NOTHING
	`
	for _, lineNgID := range lineNgIDs {
		_, err := r.db.NamedExecContext(ctx, query, lineNgID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgIDs []LineNgDiscordIDAllCoulmns) error {
	var values []string
	for _, lineNgType := range lineNgIDs {
		values = append(values, fmt.Sprintf("('%s', '%s', '%s', '%s')", lineNgType.ChannelID, lineNgType.GuildID, lineNgType.ID, lineNgType.IDType))
		_, err := r.db.ExecContext(ctx, "DELETE FROM line_ng_discord_id WHERE channel_id = $1", lineNgType.ChannelID)
		if err != nil {
			return err
		}
	}
	if len(values) == 0 {
		return nil
	}
	// INSERT されるもの以外を削除
	query := fmt.Sprintf(`
		INSERT INTO line_ng_discord_id (
			channel_id,
			guild_id,
			id,
			id_type
		) VALUES
			%s
		ON CONFLICT (channel_id, id) DO NOTHING
	`, strings.Join(values, ","))
	_, err := r.db.ExecContext(ctx, query)
	return err
}

