package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type PermissionJson struct {
	GuildID           string           `json:"guild_id"`
	PermissionCodes   []PermissionCode `json:"permission_codes"`
	PermissionUserIDs []PermissionID   `json:"permission_user_ids"`
	PermissionRoleIDs []PermissionID   `json:"permission_role_ids"`
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
	TargetID   string `json:"target_id"`
	Permission string `json:"permission"`
}
