package repository

import (
	"context"
)

type PermissionUserIDAllColumns struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	UserID     string `db:"user_id"`
	Permission string `db:"permission"`
}

type PermissionUserID struct {
	UserID     string `db:"user_id"`
	Permission string `db:"permission"`
}

func NewPermissionUserIDAllColumns(guildID, permissionType, userID, permission string) *PermissionUserIDAllColumns {
	return &PermissionUserIDAllColumns{
		GuildID:    guildID,
		Type:       permissionType,
		UserID:     userID,
		Permission: permission,
	}
}

func (r *Repository) InsertPermissionUserIDs(ctx context.Context, permissionsUserID []PermissionUserIDAllColumns) error {
	query := `
		INSERT INTO permissions_user_id (
			guild_id,
			type,
			user_id,
			permission
		) VALUES (
			:guild_id,
			:type,
			:user_id,
			:permission
		) ON CONFLICT (guild_id, type, user_id) DO NOTHING
	`
	for _, permissionUserID := range permissionsUserID {
		_, err := r.db.NamedExecContext(ctx, query, permissionUserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetGuildPermissionUserIDsAllColumns(ctx context.Context, guildID string) ([]PermissionUserIDAllColumns, error) {
	var permissionsID []PermissionUserIDAllColumns
	query := `
		SELECT
			*
		FROM
			permissions_user_id
		WHERE
			guild_id = $1
	`
	err := r.db.SelectContext(ctx, &permissionsID, query, guildID)
	return permissionsID, err
}

func (r *Repository) GetPermissionUserIDs(ctx context.Context, guildID, permissionType string) ([]PermissionUserID, error) {
	var permissionIDs []PermissionUserID
	query := `
		SELECT
			user_id,
			permission
		FROM
			permissions_user_id
		WHERE
			guild_id = $1 AND
			type = $2
	`
	err := r.db.SelectContext(ctx, &permissionIDs, query, guildID, permissionType)
	return permissionIDs, err
}

func (r *Repository) DeletePermissionUserIDs(ctx context.Context, guildId string) error {
	query := `
		DELETE FROM
			permissions_user_id
		WHERE
			guild_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, guildId)
	return err
}
