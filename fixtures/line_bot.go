package fixtures

import (
	"context"
	"testing"
)

type LineBot struct {
	GuildID          string `db:"guild_id"`
	LineNotifyToken  []byte `db:"line_notify_token"`
	LineBotToken     []byte `db:"line_bot_token"`
	LineBotSecret    []byte `db:"line_bot_secret"`
	LineGroupID      []byte `db:"line_group_id"`
	LineClientID     []byte `db:"line_client_id"`
	LineClientSecret []byte `db:"line_client_secret"`
	Iv               []byte `db:"iv"`
	DefaultChannelID string `db:"default_channel_id"`
	DebugMode        bool   `db:"debug_mode"`
}

func NewLineBot(ctx context.Context, setter ...func(b *LineBot)) *ModelConnector {
	lineBot := &LineBot{
		GuildID:          "1111111111111",
		LineNotifyToken:  []byte("1111111111111"),
		LineBotToken:     []byte("1111111111111"),
		LineBotSecret:    []byte("1111111111111"),
		LineGroupID:      []byte("1111111111111"),
		LineClientID:     []byte("1111111111111"),
		LineClientSecret: []byte("1111111111111"),
		Iv:               []byte("1111111111111"),
		DefaultChannelID: "1111111111111",
		DebugMode:        true,
	}

	return &ModelConnector{
		Model: lineBot,
		setter: func() {
			for _, s := range setter {
				s(lineBot)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineBots = append(f.LineBots, lineBot)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineBot)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_bot (
					guild_id,
					line_notify_token,
					line_bot_token,
					line_bot_secret,
					line_group_id,
					line_client_id,
					line_client_secret,
					iv,
					default_channel_id,
					debug_mode
				) VALUES (
					:guild_id,
					:line_notify_token,
					:line_bot_token,
					:line_bot_secret,
					:line_group_id,
					:line_client_id,
					:line_client_secret,
					:iv,
					:default_channel_id,
					:debug_mode
				)
			`, lineBot)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
