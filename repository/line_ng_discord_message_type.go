package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type LineNgDiscordMessageType struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	Type      int    `db:"type"`
}

func NewLineNgDiscordMessageType(channelID, guildID string, ngType int) *LineNgDiscordMessageType {
	return &LineNgDiscordMessageType{
		ChannelID: channelID,
		GuildID:   guildID,
		Type:      ngType,
	}
}

func (r *Repository) InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error {
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
	for _, lineNgType := range lineNgDiscordTypes {
		_, err := r.db.NamedExecContext(ctx, query, lineNgType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeleteMessageTypesNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordTypes []LineNgDiscordMessageType) error {
	query := `
		DELETE FROM
			line_ng_discord_message_type
		WHERE
			channel_id = ? AND
			type NOT IN (?)
	`
	if len(lineNgDiscordTypes) == 0 {
		query = `
			DELETE FROM
				line_ng_discord_message_type
			WHERE
				guild_id = $1
		`
		_, err := r.db.ExecContext(ctx, query, guildId)
		return err
	}
	typeValues := make(map[string][]int)
	for _, lineNgType := range lineNgDiscordTypes {
		typeValues[lineNgType.ChannelID] = append(typeValues[lineNgType.ChannelID], lineNgType.Type)
	}
	for channelID, types := range typeValues {
		query, args, err := db.In(query, channelID, types)
		if err != nil {
			return err
		}
		query = db.Rebind(2, query)
		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetLineNgDiscordMessageTypeByChannelID(ctx context.Context, channelID string) ([]int, error) {
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
