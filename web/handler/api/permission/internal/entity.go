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
	GuildID string `json:"guild_id"`
	Type    string `json:"type"`
	Code    int64  `json:"code"`
}

type PermissionID struct {
	GuildID    string `json:"guild_id"`
	Type       string `json:"type"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Permission string `json:"permission"`
}
