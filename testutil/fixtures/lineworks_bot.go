package fixtures

import (
	"context"
	"testing"
	"time"

	"github.com/lib/pq"
)

type LineWorksBot struct {
	GuildID               string        `db:"guild_id"`
	LineWorksBotToken     pq.ByteaArray `db:"line_works_bot_token"`
	LineWorksRefreshToken pq.ByteaArray `db:"line_works_refresh_token"`
	LineWorksGroupID      pq.ByteaArray `db:"line_works_group_id"`
	LineWorksBotID        pq.ByteaArray `db:"line_works_bot_id"`
	LineWorksBotSecret    pq.ByteaArray `db:"line_works_bot_secret"`
	RefreshTokenExpiresAt pq.NullTime   `db:"refresh_token_expires_at"`
	DefaultChannelID      string        `db:"default_channel_id"`
	DebugMode             bool          `db:"debug_mode"`
}

func NewLineWorksBot(ctx context.Context, setter ...func(b *LineWorksBot)) *ModelConnector {
	bytes := make([]byte, 13)
	copy(bytes, []byte("1111111111111"))

	lineWorksBot := &LineWorksBot{
		GuildID:               "1111111111111",
		LineWorksBotToken:     pq.ByteaArray{bytes},
		LineWorksRefreshToken: pq.ByteaArray{bytes},
		LineWorksGroupID:      pq.ByteaArray{bytes},
		LineWorksBotID:        pq.ByteaArray{bytes},
		LineWorksBotSecret:    pq.ByteaArray{bytes},
		RefreshTokenExpiresAt: pq.NullTime{Time: time.Now()},
		DefaultChannelID:      "1111111111111",
		DebugMode:             true,
	}
	return &ModelConnector{
		Model: lineWorksBot,
		setter: func() {
			for _, s := range setter {
				s(lineWorksBot)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineWorksBots = append(f.LineWorksBots, lineWorksBot)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineWorksBot)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_works_bot (
					guild_id,
					line_works_bot_token,
					line_works_refresh_token,
					line_works_group_id,
					line_works_bot_id,
					line_works_bot_secret,
					refresh_token_expires_at,
				    default_channel_id,
				    debug_mode
				) VALUES (
					:guild_id,
					:line_works_bot_token,
					:line_works_refresh_token,
					:line_works_group_id,
					:line_works_bot_id,
					:line_works_bot_secret,
					:refresh_token_expires_at,
				    :default_channel_id,
				    :debug_mode
				)
			`, lineWorksBot)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
