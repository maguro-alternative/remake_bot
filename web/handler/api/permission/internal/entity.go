package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type PermissionJson struct {
	GuildID           string             `json:"guildId"`
	PermissionCodes   []PermissionCode   `json:"permissionCodes"`
	PermissionUserIDs []PermissionUserID `json:"permissionUser_ids"`
	PermissionRoleIDs []PermissionRoleID `json:"permissionRole_ids"`
}

func (g PermissionJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.PermissionCodes, validation.Required),
	)
}

type PermissionCode struct {
	GuildID string `json:"guildId"`
	Type    string `json:"type"`
	Code    int64  `json:"code"`
}

type PermissionUserID struct {
	GuildID    string `json:"guildId"`
	Type       string `json:"type"`
	UserID     string `json:"userId"`
	Permission string `json:"permission"`
}

type PermissionRoleID struct {
	GuildID    string `json:"guildId"`
	Type       string `json:"type"`
	RoleID     string `json:"roleId"`
	Permission string `json:"permission"`
}
