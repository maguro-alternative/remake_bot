package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type LineNgDiscordRoleIDAllCoulmns struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	RoleID    string `db:"role_id"`
}

func NewLineNgDiscordRoleID(channelID, guildID, roleID string) *LineNgDiscordRoleIDAllCoulmns {
	return &LineNgDiscordRoleIDAllCoulmns{
		ChannelID: channelID,
		GuildID:   guildID,
		RoleID:    roleID,
	}
}

func (r *Repository) GetLineNgDiscordRoleIDByChannelID(ctx context.Context, channelID string) ([]string, error) {
	var ngIDs []string
	query := `
		SELECT
			role_id
		FROM
			line_ng_discord_role_id
		WHERE
			channel_id = $1
	`
	err := r.db.SelectContext(ctx, &ngIDs, query, channelID)
	return ngIDs, err
}

func (r *Repository) InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	query := `
		INSERT INTO line_ng_discord_role_id (
			channel_id,
			guild_id,
			role_id
		) VALUES (
			:channel_id,
			:guild_id,
			:role_id
		) ON CONFLICT (channel_id, role_id) DO NOTHING
	`
	for _, lineNgID := range lineNgDiscordRoleIDs {
		_, err := r.db.NamedExecContext(ctx, query, lineNgID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeleteRoleIDsNotInProvidedList(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	query := `
		DELETE FROM
			line_ng_discord_role_id
		WHERE
			channel_id = ? AND
			role_id NOT IN (?)
	`
	idValues := make(map[string][]string)
	for _, lineNgRole := range lineNgDiscordRoleIDs {
		idValues[lineNgRole.ChannelID] = append(idValues[lineNgRole.ChannelID], lineNgRole.RoleID)
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
