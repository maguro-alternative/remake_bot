package fixtures

import (
	"context"
	"testing"
)

type VcSignalNgUserID struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	UserID      string `db:"user_id"`
}

func NewVcSignalNgUserID(ctx context.Context, setter ...func(b *VcSignalNgUserID)) *ModelConnector {
	vcSignalNgUserID := &VcSignalNgUserID{
		VcChannelID: "1111111111111",
		GuildID:     "1111111111111",
		UserID:      "1111111111111",
	}

	return &ModelConnector{
		Model: vcSignalNgUserID,
		setter: func() {
			for _, s := range setter {
				s(vcSignalNgUserID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.VcSignalNgUserIDs = append(f.VcSignalNgUserIDs, vcSignalNgUserID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, vcSignalNgUserID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO vc_signal_ng_user_id (
					vc_channel_id,
					guild_id,
					user_id
				) VALUES (
					:vc_channel_id,
					:guild_id,
					:user_id
				)
			`, vcSignalNgUserID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
