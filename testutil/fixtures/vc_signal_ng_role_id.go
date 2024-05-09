package fixtures

import (
	"context"
	"testing"
)

type VcSignalNgRoleID struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	RoleID      string `db:"role_id"`
}

func NewVcSignalNgRoleID(ctx context.Context, setter ...func(b *VcSignalNgRoleID)) *ModelConnector {
	vcSignalNgRoleID := &VcSignalNgRoleID{
		VcChannelID: "1111111111111",
		GuildID:     "1111111111111",
		RoleID:      "1111111111111",
	}

	return &ModelConnector{
		Model: vcSignalNgRoleID,
		setter: func() {
			for _, s := range setter {
				s(vcSignalNgRoleID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.VcSignalNgRoleIDs = append(f.VcSignalNgRoleIDs, vcSignalNgRoleID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, vcSignalNgRoleID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO vc_signal_ng_role_id (
					vc_channel_id,
					guild_id,
					role_id
				) VALUES (
					:vc_channel_id,
					:guild_id,
					:role_id
				)
			`, vcSignalNgRoleID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
