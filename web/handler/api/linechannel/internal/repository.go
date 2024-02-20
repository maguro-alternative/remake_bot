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

func (r *Repository) UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LineChannel) error {
	query := `
		UPDATE
			line_post_discord_channel
		SET
			ng = :ng,
			bot_message = :bot_message
		WHERE
			channel_id = :channel_id
	`
	_, err := r.db.NamedExecContext(ctx, query, lineChannel)
	return err
}

func (r *Repository) InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgTypes []LineNgType) error {
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

func (r *Repository) DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgTypes []LineNgType) error {
	var values []string
	for _, lineNgType := range lineNgTypes {
		values = append(values, fmt.Sprintf("('%s', '%s', %s)", lineNgType.ChannelID, lineNgType.GuildID, lineNgType.Type))
	}
	// INSERT されるもの以外を削除
	query := fmt.Sprintf(`
	DELETE FROM
		line_ng_discord_message_type
	WHERE NOT EXISTS (
		SELECT
			*
		FROM
			(
				VALUES
					%s
			) AS t(channel_id, guild_id, type) ON CONFLICT (channel_id, type) DO NOTHING
		WHERE
			line_ng_discord_message_type.channel_id = t.channel_id AND
			line_ng_discord_message_type.type = t.type
	)
	`, strings.Join(values, ","))
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *Repository) InsertLineNgDiscordIDs(ctx context.Context, lineNgIDs []LineNgID) error {
	query := `
		INSERT INTO line_ng_id (
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

func (r *Repository) DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgIDs []LineNgID) error {
	var values []string
	for _, lineNgType := range lineNgIDs {
		values = append(values, fmt.Sprintf("('%s', '%s', '%s', '%s')", lineNgType.ChannelID, lineNgType.GuildID, lineNgType.ID, lineNgType.IDType))
	}
	// INSERT されるもの以外を削除
	query := fmt.Sprintf(`
	DELETE FROM
		line_ng_id
	WHERE NOT EXISTS (
		SELECT
			*
		FROM
			(
				VALUES
					%s
			) AS t(channel_id, guild_id, id, id_type) ON CONFLICT (channel_id, id, id_type) DO NOTHING
		WHERE
			line_ng_id.channel_id = t.channel_id AND
			line_ng_id.id = t.id AND
			line_ng_id.id_type = t.id_type
	)
	`, strings.Join(values, ","))
	_, err := r.db.ExecContext(ctx, query)
	return err
}
