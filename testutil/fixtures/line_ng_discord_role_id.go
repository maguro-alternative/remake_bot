package fixtures

import (
	"context"
	"testing"
)

type LineNgDiscordRoleID struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	RoleID    string `db:"role_id"`
}

func NewLineNgDiscordRoleID(ctx context.Context, setter ...func(b *LineNgDiscordRoleID)) *ModelConnector {
	lineNgDiscordRoleID := &LineNgDiscordRoleID{
		ChannelID: "1111111111111",
		GuildID:   "1111111111111",
		RoleID:    "1111111111111",
	}

	return &ModelConnector{
		Model: lineNgDiscordRoleID,
		setter: func() {
			for _, s := range setter {
				s(lineNgDiscordRoleID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineNgDiscordRoleIDs = append(f.LineNgDiscordRoleIDs, lineNgDiscordRoleID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineNgDiscordRoleID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_ng_discord_role_id (
					channel_id,
					guild_id,
					role_id
				) VALUES (
					:channel_id,
					:guild_id,
					:role_id
				)
			`, lineNgDiscordRoleID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
