package repository

import (
	"context"
)

type PermissionCode struct {
	GuildID string `db:"guild_id"`
	Type    string `db:"type"`
	Code    int64  `db:"code"`
}

func (r *Repository) GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error) {
	var code int64
	if permissionType == "" {
		return 8, nil
	}
	query := `
		SELECT
			code
		FROM
			permissions_code
		WHERE
			guild_id = $1 AND
			type = $2
	`
	err := r.db.GetContext(ctx, &code, query, guildID, permissionType)
	return code, err
}

func (r *Repository) GetPermissionCodes(ctx context.Context, guildID string) ([]PermissionCode, error) {
	var permissionsCode []PermissionCode
	query := `
		SELECT
			*
		FROM
			permissions_code
		WHERE
			guild_id = $1
	`
	err := r.db.SelectContext(ctx, &permissionsCode, query, guildID)
	return permissionsCode, err
}

func (r *Repository) UpdatePermissionCodes(ctx context.Context, permissionsCode []PermissionCode) error {
	query := `
		UPDATE
			permissions_code
		SET
			code = :code
		WHERE
			guild_id = :guild_id AND
			type = :type
	`
	for _, permissionCode := range permissionsCode {
		_, err := r.db.NamedExecContext(ctx, query, permissionCode)
		if err != nil {
			return err
		}
	}
	return nil
}
