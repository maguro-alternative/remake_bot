package repository

import (
	"context"
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

func (r *Repository) GetVcSignalNgRolesByChannelIDAllColumn(ctx context.Context, channelID string) ([]*VcSignalNgRoleAllColumn, error) {
	var ngRoleIDs []*VcSignalNgRoleAllColumn
	err := r.db.SelectContext(ctx, &ngRoleIDs, `
		SELECT
			*
		FROM
			vc_signal_ng_role_id
		WHERE
			vc_channel_id = $1
	`, channelID)
	if err != nil {
		return nil, err
	}
	return ngRoleIDs, nil
}

func (r *Repository) DeleteVcNgRoleByChannelID(ctx context.Context, vcChannelID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_ng_role_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
	return err
}

func (r *Repository) DeleteVcNgRoleByRoleID(ctx context.Context, roleID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_ng_role_id
		WHERE
			role_id = $1
	`, roleID)
	return err
}
