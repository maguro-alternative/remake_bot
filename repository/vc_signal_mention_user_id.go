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

func (r *Repository) GetVcSignalMentionUsersByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error) {
	var mentionUserIDs []string
	err := r.db.SelectContext(ctx, &mentionUserIDs, `
		SELECT
			user_id
		FROM
			vc_signal_mention_user_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
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

func (r *Repository) DeleteVcSignalMentionUsersByVcChannelID(ctx context.Context, vcChannelID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_mention_user_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
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

func (r *Repository) DeleteVcSignalMentionUsersNotInProvidedList(ctx context.Context, vcChannelID string, userIDs []string) error {
	query := `
		DELETE FROM
			vc_signal_mention_user_id
		WHERE
			vc_channel_id = ?
			AND user_id NOT IN (?)
	`
	if len(userIDs) == 0 {
		query = `
			DELETE FROM
				vc_signal_mention_user_id
			WHERE
				vc_channel_id = $1
		`
		_, err := r.db.ExecContext(ctx, query, vcChannelID)
		return err
	}
	query, args, err := db.In(query, vcChannelID, userIDs)
	if err != nil {
		return err
	}
	query = db.Rebind(2, query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
