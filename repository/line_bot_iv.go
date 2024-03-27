package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type LineBotIv struct {
	GuildID            string        `db:"guild_id"`
	LineNotifyTokenIv  pq.ByteaArray `db:"line_notify_token_iv"`
	LineBotTokenIv     pq.ByteaArray `db:"line_bot_token_iv"`
	LineBotSecretIv    pq.ByteaArray `db:"line_bot_secret_iv"`
	LineGroupIDIv      pq.ByteaArray `db:"line_group_id_iv"`
	LineClientIDIv     pq.ByteaArray `db:"line_client_id_iv"`
	LineClientSecretIv pq.ByteaArray `db:"line_client_secret_iv"`
}

type LineBotIvNotClient struct {
	LineNotifyTokenIv pq.ByteaArray `db:"line_notify_token_iv"`
	LineBotTokenIv    pq.ByteaArray `db:"line_bot_token_iv"`
	LineBotSecretIv   pq.ByteaArray `db:"line_bot_secret_iv"`
	LineGroupIDIv     pq.ByteaArray `db:"line_group_id_iv"`
}

func NewLineBotIv(
	guildID string,
	lineNotifyTokenIv pq.ByteaArray,
	lineBotTokenIv pq.ByteaArray,
	lineBotSecretIv pq.ByteaArray,
	lineGroupIDIv pq.ByteaArray,
	lineClientIDIv pq.ByteaArray,
	lineClientSecretIv pq.ByteaArray,
) *LineBotIv {
	return &LineBotIv{
		GuildID:            guildID,
		LineNotifyTokenIv:  lineNotifyTokenIv,
		LineBotTokenIv:     lineBotTokenIv,
		LineBotSecretIv:    lineBotSecretIv,
		LineGroupIDIv:      lineGroupIDIv,
		LineClientIDIv:     lineClientIDIv,
		LineClientSecretIv: lineClientSecretIv,
	}
}

func (r *Repository) InsertLineBotIv(ctx context.Context, guildId string) error {
	query := `
		INSERT INTO line_bot_iv (
			guild_id
		) VALUES (
			$1
		) ON CONFLICT (guild_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, guildId)
	return err
}

func (r *Repository) GetAllColumnsLineBotIv(ctx context.Context, guildID string) (LineBotIv, error) {
	var lineBotIv LineBotIv
	query := `
		SELECT
			line_notify_token_iv,
			line_bot_token_iv,
			line_bot_secret_iv,
			line_client_id_iv,
			line_client_secret_iv,
			line_group_id_iv
		FROM
			line_bot_iv
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineBotIv, query, guildID)
	return lineBotIv, err
}

func (r *Repository) GetLineBotIvNotClient(ctx context.Context, guildID string) (LineBotIvNotClient, error) {
	var lineBotIv LineBotIvNotClient
	query := `
		SELECT
			line_notify_token_iv,
			line_bot_token_iv,
			line_bot_secret_iv,
			line_group_id_iv
		FROM
			line_bot_iv
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineBotIv, query, guildID)
	return lineBotIv, err
}


func (r *Repository) UpdateLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineBotIv.LineNotifyTokenIv) > 0 && len(lineBotIv.LineNotifyTokenIv[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_notify_token_iv = :line_notify_token_iv")
	}
	if len(lineBotIv.LineBotTokenIv) > 0 && len(lineBotIv.LineBotTokenIv[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_token_iv = :line_bot_token_iv")
	}
	if len(lineBotIv.LineBotSecretIv) > 0 && len(lineBotIv.LineBotSecretIv[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_bot_secret_iv = :line_bot_secret_iv")
	}
	if len(lineBotIv.LineGroupIDIv) > 0 && len(lineBotIv.LineGroupIDIv[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_group_id_iv = :line_group_id_iv")
	}
	if len(lineBotIv.LineClientIDIv) > 0 && len(lineBotIv.LineClientIDIv[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_client_id_iv = :line_client_id_iv")
	}
	if len(lineBotIv.LineClientSecretIv) > 0 && len(lineBotIv.LineClientSecretIv[0]) > 0 {
		setQueryArray = append(setQueryArray, "line_client_secret_iv = :line_client_secret_iv")
	}
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE
			line_bot_iv
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)
	_, err := r.db.NamedExecContext(ctx, query, lineBotIv)
	return err
}
