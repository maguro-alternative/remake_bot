package repository

import (
	"context"
)

type PermissionRoleIDAllColumns struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	RoleID     string `db:"role_id"`
	Permission string `db:"permission"`
}

type PermissionRoleID struct {
	RoleID     string `db:"target_id"`
	Permission string `db:"permission"`
}

func NewPermissionRoleIDAllColumns(guildID, permissionType, targetType, roleID, permission string) *PermissionRoleIDAllColumns {
	return &PermissionRoleIDAllColumns{
		GuildID:    guildID,
		Type:       permissionType,
		RoleID:     roleID,
		Permission: permission,
	}
}

func (r *Repository) InsertPermissionRoleIDs(ctx context.Context, permissionsRoleID []PermissionRoleIDAllColumns) error {
	query := `
		INSERT INTO permissions_role_id (
			guild_id,
			type,
			role_id,
			permission
		) VALUES (
			:guild_id,
			:type,
			:role_id,
			:permission
		) ON CONFLICT (guild_id, type, role_id) DO NOTHING
	`
	for _, permissionRoleID := range permissionsRoleID {
		_, err := r.db.NamedExecContext(ctx, query, permissionRoleID)
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
			role_id,
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
