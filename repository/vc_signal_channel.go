package repository

import (
	"context"
)

type VcSignalChannelAllColumns struct {
	VcChannelID     string `db:"vc_channel_id"`
	GuildID         string `db:"guild_id"`
	SendSignal      bool   `db:"send_signal"`
	SendChannelID   string `db:"send_channel_id"`
	JoinBot         bool   `db:"join_bot"`
	EveryoneMention bool   `db:"everyone_mention"`
}

type VcSignalChannelNotGuildID struct {
	VcChannelID     string `db:"vc_channel_id"`
	GuildID         string `db:"guild_id"`
	SendSignal      bool   `db:"send_signal"`
	SendChannelID   string `db:"send_channel_id"`
	JoinBot         bool   `db:"join_bot"`
	EveryoneMention bool   `db:"everyone_mention"`
}

func (r *Repository) InsertVcSignalChannel(ctx context.Context, vcChannelID string, guildID, sendChannelID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO vc_signal_channel (
			vc_channel_id,
			guild_id,
			send_signal,
			send_channel_id,
			join_bot,
			everyone_mention
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		) ON CONFLICT (vc_channel_id) DO NOTHING
	`, vcChannelID, guildID, true, sendChannelID, false, true)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetVcSignalChennel(ctx context.Context, vcChannelID string) (*VcSignalChannelAllColumns, error) {
	var vcSignalChannel VcSignalChannelAllColumns
	err := r.db.GetContext(ctx, &vcSignalChannel, "SELECT * FROM vc_signal_channel WHERE vc_channel_id = ?", vcChannelID)
	if err != nil {
		return nil, err
	}
	return &vcSignalChannel, nil
}

func (r *Repository) UpdateVcSignalChannel(ctx context.Context, vcChannel VcSignalChannelNotGuildID) error {
	_, err := r.db.ExecContext(ctx, `
	UPDATE
		vc_signal_channel
	SET
		send_signal = :send_signal,
		send_channel_id = :send_channel_id,
		join_bot = :join_bot,
		everyone_mention = :everyone_mention
	WHERE
		vc_channel_id = :vc_channel_id
	`, vcChannel)
	return err
}

func (r *Repository) DeleteVcSignalChannel(ctx context.Context, vcChannelID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM vc_signal_channel WHERE vc_channel_id = ?", vcChannelID)
	if err != nil {
		return err
	}
	return nil
}
