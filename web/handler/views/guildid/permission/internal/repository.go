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

func (r *Repository) GetPermissionCode(ctx context.Context, guildID string) (PermissionCode, error) {
	var permissionCode PermissionCode
	query := `
		SELECT
			*
		FROM
			permissions_code
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &permissionCode, query, guildID)
	return permissionCode, err
}

func (r *Repository) GetPermissionID(ctx context.Context, guildID string) (PermissionID, error) {
	var permissionID PermissionID
	query := `
		SELECT
			*
		FROM
			permissions_id
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &permissionID, query, guildID)
	return permissionID, err
}