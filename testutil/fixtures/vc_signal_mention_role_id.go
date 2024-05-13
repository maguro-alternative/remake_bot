package fixtures

import (
	"context"
	"testing"
)

type VcSignalMentionRoleID struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	RoleID      string `db:"role_id"`
}

func NewVcSignalMentionRoleID(ctx context.Context, setter ...func(b *VcSignalMentionRoleID)) *ModelConnector {
	vcSignalMentionRoleID := &VcSignalMentionRoleID{
		VcChannelID: "1111111111111",
		GuildID:     "1111111111111",
		RoleID:      "1111111111111",
	}

	return &ModelConnector{
		Model: vcSignalMentionRoleID,
		setter: func() {
			for _, s := range setter {
				s(vcSignalMentionRoleID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.VcSignalMentionRoleIDs = append(f.VcSignalMentionRoleIDs, vcSignalMentionRoleID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, vcSignalMentionRoleID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO vc_signal_mention_role_id (
					vc_channel_id,
					guild_id,
					role_id
				) VALUES (
					:vc_channel_id,
					:guild_id,
					:role_id
				)
			`, vcSignalMentionRoleID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
