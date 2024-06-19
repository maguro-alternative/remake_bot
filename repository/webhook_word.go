package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type WebhookWord struct {
	ID        int64  `db:"id"`
	Condition string `db:"condition"`
	Word      string `db:"word"`
}

func (r *Repository) InsertWebhookWord(
	ctx context.Context,
	condition string,
	word string,
) error {
	query := `
		INSERT INTO webhook_word (
			condition,
			word
		) VALUES (
			$1,
			$2
		) ON CONFLICT (condition, word) DO NOTHING
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
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
			condition = $2
	`
	var webhookWord []*WebhookWord
	err := r.db.SelectContext(ctx, &webhookWord, query, webhookSerialID, condition)
	return webhookWord, err
}

func (r *Repository) DeleteWebhookWordWithWebhookSerialIDAndCondition(
	ctx context.Context,
	webhookSerialID int64,
	condition string,
) error {
	query := `
		DELETE FROM
			webhook_word
		WHERE
			webhook_serial_id = $1
			condition = $2
	`
	_, err := r.db.ExecContext(ctx, query, webhookSerialID, condition)
	return err
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
			condition = ? AND
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
