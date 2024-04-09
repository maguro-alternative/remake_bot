package fixtures

import (
	"context"
	"testing"
)

type VcSignalNgID struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	IDType      string `db:"id_type"`
	ID          string `db:"id"`
}

func NewVcSignalNgID(ctx context.Context, setter ...func(b *VcSignalNgID)) *ModelConnector {
	vcSignalNgID := &VcSignalNgID{
		VcChannelID: "1111111111111",
		GuildID:     "1111111111111",
		IDType:      "user",
		ID:          "1111111111111",
	}

	return &ModelConnector{
		Model: vcSignalNgID,
		setter: func() {
			for _, s := range setter {
				s(vcSignalNgID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.VcSignalNgIDs = append(f.VcSignalNgIDs, vcSignalNgID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, vcSignalNgID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO vc_signal_ng_id (
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
			`, vcSignalNgID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
