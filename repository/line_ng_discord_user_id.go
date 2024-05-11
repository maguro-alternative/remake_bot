package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
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

func (r *Repository) DeleteUserIDsNotInProvidedList(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error {
	query := `
		DELETE FROM
			line_ng_discord_user_id
		WHERE
			channel_id = ? AND
			user_id NOT IN (?)
	`
	idValues := make(map[string][]string)
	for _, lineNgUser := range lineNgDiscordUserIDs {
		idValues[lineNgUser.ChannelID] = append(idValues[lineNgUser.ChannelID], lineNgUser.UserID)
	}
	for channelID, roleIDs := range idValues {
		query, args, err := db.In(query, channelID, roleIDs)
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
