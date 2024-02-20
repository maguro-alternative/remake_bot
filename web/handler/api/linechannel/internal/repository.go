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

func (r *Repository) UpdateLineChannel(ctx context.Context, lineChannel []LineChannel) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineChannel) > 0 {
		setQueryArray = append(setQueryArray, "ng = :ng")
	}
	if len(lineChannel) > 0 {
		setQueryArray = append(setQueryArray, "bot_message = :bot_message")
	}
	setNameQuery = strings.Join(setQueryArray, ",")

	query := fmt.Sprintf(`
		UPDATE
			line_channel
		SET
			%s
		WHERE
			channel_id = :channel_id
	`, setNameQuery)
	_, err := r.db.NamedExecContext(ctx, query, lineChannel)
	return err
}

func (r *Repository) InsertLineNgTypes(ctx context.Context, lineNgTypes []LineNgType) error {
	query := `
		INSERT INTO line_ng_type (
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

func (r *Repository) DeleteNotInsertLineNgTypes(ctx context.Context, lineNgTypes []LineNgType) error {
	var values []string
	for _, lineNgType := range lineNgTypes {
		values = append(values, fmt.Sprintf("('%s', '%s', %s)", lineNgType.ChannelID, lineNgType.GuildID, lineNgType.Type))
	}
	// INSERT されるもの以外を削除
	query := fmt.Sprintf(`
	DELETE FROM
		line_ng_type
	WHERE NOT EXISTS (
		SELECT
			*
		FROM
			(
				VALUES
					%s
			) AS t(channel_id, guild_id, type) ON CONFLICT (channel_id, type) DO NOTHING
		WHERE
			line_ng_type.channel_id = t.channel_id AND
			line_ng_type.type = t.type
	)
	`, strings.Join(values, ","))
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *Repository) InsertLineNgIDs(ctx context.Context, lineNgIDs []LineNgID) error {
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

func (r *Repository) DeleteNotInsertLineNgIDs(ctx context.Context, lineNgIDs []LineNgID) error {
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
