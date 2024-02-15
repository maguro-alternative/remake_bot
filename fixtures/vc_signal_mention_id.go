package fixtures

import (
	"context"
	"testing"
)

type VcSignalMentionID struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	IDType      string `db:"id_type"`
	ID          string `db:"id"`
}

func NewVcSignalMentionID(ctx context.Context, setter ...func(b *VcSignalMentionID)) *ModelConnector {
	vcSignalMentionID := &VcSignalMentionID{
		VcChannelID: "1111111111111",
		GuildID:     "1111111111111",
		IDType:      "user",
		ID:          "1111111111111",
	}

	return &ModelConnector{
		Model: vcSignalMentionID,
		setter: func() {
			for _, s := range setter {
				s(vcSignalMentionID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.VcSignalMentionIDs = append(f.VcSignalMentionIDs, vcSignalMentionID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, vcSignalMentionID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO vc_signal_mention_id (
					vc_channel_id,
					guild_id,
					id_type,
					id
				) VALUES (
					:vc_channel_id,
					:guild_id,
					:id_type,
					:id
				)
			`, vcSignalMentionID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
