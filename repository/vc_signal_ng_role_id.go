package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type VcSignalNgRoleAllColumn struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	RoleID      string `db:"role_id"`
}

func (r *Repository) InsertVcSignalNgRole(ctx context.Context, vcChannelID, guildID, roleID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO vc_signal_ng_role_id (
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

func (r *Repository) GetVcSignalNgRolesByVcChannelIDAllColumn(ctx context.Context, vcChannelID string) ([]*VcSignalNgRoleAllColumn, error) {
	var ngRoleIDs []*VcSignalNgRoleAllColumn
	err := r.db.SelectContext(ctx, &ngRoleIDs, `
		SELECT
			*
		FROM
			vc_signal_ng_role_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
	if err != nil {
		return nil, err
	}
	return ngRoleIDs, nil
}

func (r *Repository) DeleteVcSignalNgRoleByVcChannelID(ctx context.Context, vcChannelID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_ng_role_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
	return err
}

func (r *Repository) DeleteVcSignalNgRoleByGuildID(ctx context.Context, guildID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_ng_role_id
		WHERE
			guild_id = $1
	`, guildID)
	return err
}

func (r *Repository) DeleteVcSignalNgRoleByRoleID(ctx context.Context, roleID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_ng_role_id
		WHERE
			role_id = $1
	`, roleID)
	return err
}

func (r *Repository) DeleteVcSignalRolesNotInProvidedList(ctx context.Context, vcChannelID string, roleIDs []string) error {
	query := `
		DELETE FROM
			vc_signal_ng_role_id
		WHERE
			vc_channel_id = ? AND
			role_id NOT IN (?)
	`
	if len(roleIDs) == 0 {
		query = `
			DELETE FROM
				vc_signal_ng_role_id
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