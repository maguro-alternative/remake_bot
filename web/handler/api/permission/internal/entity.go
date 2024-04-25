package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type PermissionJson struct {
	GuildID           string             `json:"guild_id"`
	PermissionCodes   []PermissionCode   `json:"permission_codes"`
	PermissionUserIDs []PermissionUserID `json:"permission_user_ids"`
	PermissionRoleIDs []PermissionRoleID `json:"permission_role_ids"`
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

type PermissionUserID struct {
	GuildID    string `json:"guild_id"`
	Type       string `json:"type"`
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
}

type PermissionRoleID struct {
	GuildID    string `json:"guild_id"`
	Type       string `json:"type"`
	RoleID     string `json:"role_id"`
	Permission string `json:"permission"`
}
