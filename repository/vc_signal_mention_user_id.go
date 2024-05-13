package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type VcSignalMentionUser struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	UserID      string `db:"user_id"`
}

func (r *Repository) InsertVcSignalMentionUser(ctx context.Context, vcChannelID, guildID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO vc_signal_mention_user_id (
			vc_channel_id,
			guild_id,
			user_id
		) VALUES (
			$1,
			$2,
			$3
		) ON CONFLICT (vc_channel_id, user_id) DO NOTHING
	`, vcChannelID, guildID, userID)
	return err
}

func (r *Repository) GetVcSignalMentionUsersByChannelID(ctx context.Context, channelID string) ([]*VcSignalMentionUser, error) {
	var mentionUserIDs []*VcSignalMentionUser
	err := r.db.SelectContext(ctx, &mentionUserIDs, `
		SELECT
			*
		FROM
			vc_signal_mention_user_id
		WHERE
			vc_channel_id = $1
	`, channelID)
	if err != nil {
		return nil, err
	}
	return mentionUserIDs, nil
}

func (r *Repository) DeleteVcSignalMentionUser(ctx context.Context, vcChannelID, guildID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_user_id
		WHERE
			vc_channel_id = $1
			AND guild_id = $2
			AND user_id = $3
	`, vcChannelID, guildID, userID)
	return err
}

func (r *Repository) DeleteVcSignalMentionUsersByChannelID(ctx context.Context, channelID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_user_id
		WHERE
			vc_channel_id = $1
	`, channelID)
	return err
}

func (r *Repository) DeleteVcSignalMentionUsersByGuildID(ctx context.Context, guildID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_user_id
		WHERE
			guild_id = $1
	`, guildID)
	return err
}

func (r *Repository) DeleteVcSignalMentionUsersByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_user_id
		WHERE
			user_id = $1
	`, userID)
	return err
}

func (r *Repository) DeleteVcSignalMentionUsersNotInProvidedList(ctx context.Context, channelID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	query := `
		DELETE FROM
			vc_signal_mention_user_id
		WHERE
			vc_channel_id = ?
			AND user_id NOT IN (?)
	`
	query, args, err := db.In(query, channelID, userIDs)
	if err != nil {
		return err
	}
	query = db.Rebind(2, query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
