package fixtures

import (
	"context"
	"testing"
)

type LineChannel struct {
	ChannelID  string `db:"channel_id"`
	GuildID    string `db:"guild_id"`
	Ng         bool   `db:"ng"`
	BotMessage bool   `db:"bot_message"`
}

func NewLineChannel(ctx context.Context, setter ...func(b *LineChannel)) *ModelConnector {
	lineChannel := &LineChannel{
		ChannelID:  "1111111111111",
		GuildID:    "1111111111111",
		Ng:         false,
		BotMessage: false,
	}

	return &ModelConnector{
		Model: lineChannel,
		setter: func() {
			for _, s := range setter {
				s(lineChannel)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineChannels = append(f.LineChannels, lineChannel)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineChannel)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_post_discord_channel (
					channel_id,
					guild_id,
					ng,
					bot_message
				) VALUES (
					:channel_id,
					:guild_id,
					:ng,
					:bot_message
				)
			`, lineChannel)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
