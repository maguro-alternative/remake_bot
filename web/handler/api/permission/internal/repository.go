package internal

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type Repository struct {
	db db.Driver
}

func NewRepository(db db.Driver) *Repository {
	return &Repository{
		db: db,
	}
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

func (r *Repository) InsertPermissionIDs(ctx context.Context, permissionsID []PermissionID) error {
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
