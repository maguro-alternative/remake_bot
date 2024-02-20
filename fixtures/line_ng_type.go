package fixtures

import (
	"context"
	"testing"
)

type LineNgType struct {
	ChannelID string `db:"channel_id"`
	GuildID   string `db:"guild_id"`
	Type      int    `db:"type"`
}

func NewLineNgType(ctx context.Context, setter ...func(b *LineNgType)) *ModelConnector {
	lineNgType := &LineNgType{
		ChannelID: "1111111111111",
		GuildID:   "2222222222222",
		Type:      6,
	}

	return &ModelConnector{
		Model: lineNgType,
		setter: func() {
			for _, s := range setter {
				s(lineNgType)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineNgTypes = append(f.LineNgTypes, lineNgType)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineNgType)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_ng_discord_message_type (
					channel_id,
					guild_id,
					type
				) VALUES (
					:channel_id,
					:guild_id,
					:type
				)
			`, lineNgType)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
