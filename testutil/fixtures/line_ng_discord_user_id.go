package fixtures

import (
	"context"
	"testing"
)

type LineNgDiscordUserID struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	ID        string `db:"id"`
}

func NewLineNgDiscordUserID(ctx context.Context, setter ...func(b *LineNgDiscordUserID)) *ModelConnector {
	lineNgDiscordUserID := &LineNgDiscordUserID{
		ChannelID: "1111111111111",
		GuildID:   "1111111111111",
		ID:        "1111111111111",
	}

	return &ModelConnector{
		Model: lineNgDiscordUserID,
		setter: func() {
			for _, s := range setter {
				s(lineNgDiscordUserID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineNgDiscordUserIDs = append(f.LineNgDiscordUserIDs, lineNgDiscordUserID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineNgDiscordUserID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_ng_discord_user_id (
					channel_id,
					guild_id,
					id
				) VALUES (
					:channel_id,
					:guild_id,
					:id
				)
			`, lineNgDiscordUserID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
