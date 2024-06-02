package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type VcSignalMentionRole struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	RoleID      string `db:"role_id"`
}

func (r *Repository) InsertVcSignalMentionRole(ctx context.Context, vcChannelID, guildID, roleID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO vc_signal_mention_role_id (
			vc_channel_id,
			guild_id,
			role_id
		) VALUES (
			$1,
			$2,
			$3
		) ON CONFLICT (vc_channel_id, role_id) DO NOTHING
	`, vcChannelID, guildID, roleID)
	return err
}

func (r *Repository) GetVcSignalMentionRolesByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error) {
	var mentionRoleIDs []string
	err := r.db.SelectContext(ctx, &mentionRoleIDs, `
		SELECT
			role_id
		FROM
			vc_signal_mention_role_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
	if err != nil {
		return nil, err
	}
	return mentionRoleIDs, nil
}

func (r *Repository) DeleteVcSignalMentionRole(ctx context.Context, vcChannelID, guildID, roleID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_role_id
		WHERE
			vc_channel_id = $1
			AND guild_id = $2
			AND role_id = $3
	`, vcChannelID, guildID, roleID)
	return err
}

func (r *Repository) DeleteVcSignalMentionRolesByVcChannelID(ctx context.Context, vcChannelID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_role_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
	return err
}

func (r *Repository) DeleteVcSignalMentionRolesByGuildID(ctx context.Context, guildID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_role_id
		WHERE
			guild_id = $1
	`, guildID)
	return err
}

func (r *Repository) DeleteVcSignalMentionRolesByRoleID(ctx context.Context, roleID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_role_id
		WHERE
			role_id = $1
	`, roleID)
	return err
}

func (r *Repository) DeleteVcSignalMentionRolesNotInProvidedList(ctx context.Context, vcChannelID string, roleIDs []string) error {
	query := `
		DELETE FROM
			vc_signal_mention_role_id
		WHERE
			vc_channel_id = ?
			AND role_id NOT IN (?)
	`
	if len(roleIDs) == 0 {
		query = `
			DELETE FROM
				vc_signal_mention_role_id
			WHERE
				vc_channel_id = $1
		`
		_, err := r.db.ExecContext(ctx, query, vcChannelID)
		return err
	}
	query, args, err := db.In(query, vcChannelID, roleIDs)
	if err != nil {
		return err
	}
	query = db.Rebind(2, query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
