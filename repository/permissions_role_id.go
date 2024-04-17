package repository

import (
	"context"
)

type PermissionRoleIDAllColumns struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	TargetID   string `db:"target_id"`
	Permission string `db:"permission"`
}

type PermissionRoleID struct {
	TargetID   string `db:"target_id"`
	Permission string `db:"permission"`
}

func NewPermissionRoleIDAllColumns(guildID, permissionType, targetType, targetID, permission string) *PermissionRoleIDAllColumns {
	return &PermissionRoleIDAllColumns{
		GuildID:    guildID,
		Type:       permissionType,
		TargetID:   targetID,
		Permission: permission,
	}
}

func (r *Repository) InsertPermissionRoleIDs(ctx context.Context, permissionsID []PermissionRoleIDAllColumns) error {
	query := `
		INSERT INTO permissions_role_id (
			guild_id,
			type,
			target_id,
			permission
		) VALUES (
			:guild_id,
			:type,
			:target_id,
			:permission
		) ON CONFLICT (guild_id, type, target_id) DO NOTHING
	`
	for _, permissionID := range permissionsID {
		_, err := r.db.NamedExecContext(ctx, query, permissionID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetGuildPermissionRoleIDsAllColumns(ctx context.Context, guildID string) ([]PermissionRoleIDAllColumns, error) {
	var permissionsID []PermissionRoleIDAllColumns
	query := `
		SELECT
			*
		FROM
			permissions_role_id
		WHERE
			guild_id = $1
	`
	err := r.db.SelectContext(ctx, &permissionsID, query, guildID)
	return permissionsID, err
}

func (r *Repository) GetPermissionRoleIDs(ctx context.Context, guildID, permissionType string) ([]PermissionRoleID, error) {
	var permissionIDs []PermissionRoleID
	query := `
		SELECT
			target_id,
			permission
		FROM
			permissions_role_id
		WHERE
			guild_id = $1 AND
			type = $2
	`
	err := r.db.SelectContext(ctx, &permissionIDs, query, guildID, permissionType)
	return permissionIDs, err
}

func (r *Repository) DeletePermissionRoleIDs(ctx context.Context, guildId string) error {
	query := `
		DELETE FROM
			permissions_role_id
		WHERE
			guild_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, guildId)
	return err
}
