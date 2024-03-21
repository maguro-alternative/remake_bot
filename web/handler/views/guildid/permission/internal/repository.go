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

func (r *Repository) GetPermissionIDs(ctx context.Context, guildID string) ([]PermissionID, error) {
	var permissionsID []PermissionID
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