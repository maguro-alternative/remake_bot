package repository

import (
	"context"
)

type PermissionIDAllColumns struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	TargetType string `db:"target_type"`
	TargetID   string `db:"target_id"`
	Permission string `db:"permission"`
}

type PermissionID struct {
	TargetType string `db:"target_type"`
	TargetID   string `db:"target_id"`
	Permission string `db:"permission"`
}

func (r *Repository) InsertPermissionIDs(ctx context.Context, permissionsID []PermissionIDAllColumns) error {
	query := `
		INSERT INTO permissions_id (
			guild_id,
			type,
			target_type,
			target_id,
			permission
		) VALUES (
			:guild_id,
			:type,
			:target_type,
			:target_id,
			:permission
		) ON CONFLICT (guild_id, type, target_type, target_id) DO NOTHING
	`
	for _, permissionID := range permissionsID {
		_, err := r.db.NamedExecContext(ctx, query, permissionID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetGuildPermissionIDsAllColumns(ctx context.Context, guildID string) ([]PermissionIDAllColumns, error) {
	var permissionsID []PermissionIDAllColumns
	query := `
		SELECT
			*
		FROM
			permissions_id
		WHERE
			guild_id = $1
	`
	err := r.db.SelectContext(ctx, &permissionsID, query, guildID)
	return permissionsID, err
}

func (r *Repository) GetPermissionIDs(ctx context.Context, guildID, permissionType string) ([]PermissionID, error) {
	var permissionIDs []PermissionID
	query := `
		SELECT
			target_type,
			target_id,
			permission
		FROM
			permissions_id
		WHERE
			guild_id = $1 AND
			type = $2
	`
	err := r.db.SelectContext(ctx, &permissionIDs, query, guildID, permissionType)
	return permissionIDs, err
}

func (r *Repository) DeletePermissionIDs(ctx context.Context, guildId string) error {
	query := `
		DELETE FROM
			permissions_id
		WHERE
			guild_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, guildId)
	return err
}
