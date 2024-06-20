package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type WebhookWord struct {
	WebhookSerialID int64  `db:"webhook_serial_id"`
	Condition       string `db:"conditions"`
	Word            string `db:"word"`
}

func (r *Repository) InsertWebhookWord(
	ctx context.Context,
	webhookSerialID int64,
	condition string,
	word string,
) error {
	query := `
		INSERT INTO webhook_word (
			webhook_serial_id,
			conditions,
			word
		) VALUES (
			$1,
			$2,
			$3
		) ON CONFLICT (webhook_serial_id, word) DO NOTHING
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		webhookSerialID,
		condition,
		word,
	)
	return err
}

func (r *Repository) GetWebhookWordWithWebhookSerialIDAndCondition(
	ctx context.Context,
	webhookSerialID int64,
	condition string,
) ([]*WebhookWord, error) {
	query := `
		SELECT
			*
		FROM
			webhook_word
		WHERE
			webhook_serial_id = $1
			conditions = $2
	`
	var webhookWord []*WebhookWord
	err := r.db.SelectContext(ctx, &webhookWord, query, webhookSerialID, condition)
	return webhookWord, err
}

func (r *Repository) DeleteWebhookWordsNotInProvidedList(
	ctx context.Context,
	webhookSerialID int64,
	conditions string,
	words []string,
) error {
	query := `
		DELETE FROM
			webhook_word
		WHERE
			webhook_serial_id = ? AND
			conditions = ? AND
			word NOT IN (?)
	`
	if len(words) == 0 {
		query = `
			DELETE FROM
				webhook_word
			WHERE
				webhook_serial_id = $1 AND
				condition = $2
		`
		_, err := r.db.ExecContext(ctx, query, webhookSerialID, conditions)
		return err
	}
	query, args, err := db.In(query, webhookSerialID, conditions, words)
	if err != nil {
		return err
	}
	query = db.Rebind(2, query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
