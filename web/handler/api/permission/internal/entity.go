package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type PermissionJson struct {
	GuildID         string           `json:"guild_id"`
	PermissionCodes []PermissionCode `json:"permission_codes"`
	PermissionIDs   []PermissionID   `json:"permission_ids"`
}

func (g PermissionJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.PermissionCodes, validation.Required),
	)
}

type PermissionCode struct {
	GuildID string `db:"guild_id" json:"guild_id"`
	Type    string `db:"type" json:"type"`
	Code    int64  `db:"code" json:"code"`
}

type PermissionID struct {
	GuildID    string `db:"guild_id" json:"guild_id"`
	Type       string `db:"type" json:"type"`
	TargetType string `db:"target_type" json:"target_type"`
	TargetID   string `db:"target_id" json:"target_id"`
	Permission string `db:"permission" json:"permission"`
}
