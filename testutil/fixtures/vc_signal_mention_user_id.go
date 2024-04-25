package fixtures

import (
	"context"
	"testing"
)

type VcSignalMentionUserID struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	UserID      string `db:"user_id"`
}

func NewVcSignalMentionUserID(ctx context.Context, setter ...func(b *VcSignalMentionUserID)) *ModelConnector {
	vcSignalMentionUserID := &VcSignalMentionUserID{
		VcChannelID: "1111111111111",
		GuildID:     "1111111111111",
		UserID:      "1111111111111",
	}

	return &ModelConnector{
		Model: vcSignalMentionUserID,
		setter: func() {
			for _, s := range setter {
				s(vcSignalMentionUserID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.VcSignalMentionUserIDs = append(f.VcSignalMentionUserIDs, vcSignalMentionUserID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, vcSignalMentionUserID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO vc_signal_mention_user_id (
					vc_channel_id,
					guild_id,
					user_id
				) VALUES (
					:vc_channel_id,
					:guild_id,
					:id_type,
					:user_id
				)
			`, vcSignalMentionUserID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
