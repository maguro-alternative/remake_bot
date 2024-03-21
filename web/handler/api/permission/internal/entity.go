package internal

type PermissionJson struct {
	GuildID string `json:"guild_id"`
	PermissionCodes []struct {
		Type    string `json:"type"`
		Code    int64  `json:"code"`
	} `json:"permission_codes"`
	PermissionIDs []struct {
		Type       string `json:"type"`
		TargetType string `json:"target_type"`
		TargetID   string `json:"target_id"`
		Permission string `json:"permission"`
	} `json:"permission_ids"`
}

type PermissionCode struct {
	GuildID string `db:"guild_id"`
	Type    string `db:"type"`
	Code    int64  `db:"code"`
}

type PermissionID struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	TargetType string `db:"target_type"`
	TargetID   string `db:"target_id"`
	Permission string `db:"permission"`
}
