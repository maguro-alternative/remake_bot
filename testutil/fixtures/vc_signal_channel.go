package fixtures

import (
	"context"
	"testing"
)

type VcSignalChannel struct {
	VcChannelID     string `db:"vc_channel_id"`
	GuildID         string `db:"guild_id"`
	SendSignal      bool   `db:"send_signal"`
	SendChannelID   string `db:"send_channel_id"`
	JoinBot         bool   `db:"join_bot"`
	EveryOneMention bool   `db:"everyone_mention"`
}

func NewVcSignalChannel(ctx context.Context, setter ...func(b *VcSignalChannel)) *ModelConnector {
	vcSignalChannel := &VcSignalChannel{
		VcChannelID:     "1111111111111",
		GuildID:         "1111111111111",
		SendSignal:      true,
		SendChannelID:   "1111111111111",
		JoinBot:         false,
		EveryOneMention: true,
	}

	return &ModelConnector{
		Model: vcSignalChannel,
		setter: func() {
			for _, s := range setter {
				s(vcSignalChannel)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.VcSignalChannels = append(f.VcSignalChannels, vcSignalChannel)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, vcSignalChannel)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO vc_signal_channel (
					vc_channel_id,
					guild_id,
					send_signal,
					send_channel_id,
					join_bot,
					everyone_mention
				) VALUES (
					:vc_channel_id,
					:guild_id,
					:send_signal,
					:send_channel_id,
					:join_bot,
					:everyone_mention
				)
			`, vcSignalChannel)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
