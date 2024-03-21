package internal

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
