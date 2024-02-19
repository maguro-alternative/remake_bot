package fixtures

import (
	"context"
	"testing"
)

type LineNgDiscordID struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	IDType    string `db:"id_type"`
	ID        string `db:"id"`
}

func NewLineNgDiscordID(ctx context.Context, setter ...func(b *LineNgDiscordID)) *ModelConnector {
	lineNgDiscordID := &LineNgDiscordID{
		ChannelID: "1111111111111",
		GuildID:   "1111111111111",
		IDType:    "user",
		ID:        "1111111111111",
	}

	return &ModelConnector{
		Model: lineNgDiscordID,
		setter: func() {
			for _, s := range setter {
				s(lineNgDiscordID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineNgDiscordIDs = append(f.LineNgDiscordIDs, lineNgDiscordID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineNgDiscordID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_ng_discord_id (
					channel_id,
					guild_id,
					id_type,
					id
				) VALUES (
					:channel_id,
					:guild_id,
					:id_type,
					:id
				)
			`, lineNgDiscordID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
