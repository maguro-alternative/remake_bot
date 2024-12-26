package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type LineWorksBotInfoIV struct {
	GuildID string `db:"guild_id"`
	LineWorksClientIDIV pq.ByteaArray `db:"line_works_client_id_iv"`
	LineWorksClientSecretIV pq.ByteaArray `db:"line_works_client_secret_iv"`
	LineWorksServiceAccountIV pq.ByteaArray `db:"line_works_service_account_iv"`
	LineWorksPrivateKeyIV pq.ByteaArray `db:"line_works_private_key_iv"`
	LineWorksDomainIDIV pq.ByteaArray `db:"line_works_domain_id_iv"`
	LineWorksAdminIDIV pq.ByteaArray `db:"line_works_admin_id_iv"`
}

func NewLineWorksBotInfoIV(
	guildID string,
	lineWorksClientIDIV pq.ByteaArray,
	lineWorksClientSecretIV pq.ByteaArray,
	lineWorksServiceAccountIV pq.ByteaArray,
	lineWorksPrivateKeyIV pq.ByteaArray,
	lineWorksDomainIDIV pq.ByteaArray,
	lineWorksAdminIDIV pq.ByteaArray,
) *LineWorksBotInfoIV {
	return &LineWorksBotInfoIV{
		GuildID: guildID,
		LineWorksClientIDIV: lineWorksClientIDIV,
		LineWorksClientSecretIV: lineWorksClientSecretIV,
		LineWorksServiceAccountIV: lineWorksServiceAccountIV,
		LineWorksPrivateKeyIV: lineWorksPrivateKeyIV,
		LineWorksDomainIDIV: lineWorksDomainIDIV,
		LineWorksAdminIDIV: lineWorksAdminIDIV,
	}
}

func (r *Repository) InsertLineWorksBotInfoIV(ctx context.Context, lineWorksBotInfoIv *LineWorksBotInfoIV) error {
	query := `
		INSERT INTO
			line_works_bot_info_iv (
				guild_id,
				line_works_client_id_iv,
				line_works_client_secret_iv,
				line_works_service_account_iv,
				line_works_private_key_iv,
				line_works_domain_id_iv,
				line_works_admin_id_iv
			)
		VALUES (
			:guild_id,
			:line_works_client_id_iv,
			:line_works_client_secret_iv,
			:line_works_service_account_iv,
			:line_works_private_key_iv,
			:line_works_domain_id_iv,
			:line_works_admin_id_iv
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, lineWorksBotInfoIv)
	return err
}

func (r *Repository) GetLineWorksBotInfoIVByGuildID(ctx context.Context, guildID string) (*LineWorksBotInfoIV, error) {
	var lineWorksBotInfoIV LineWorksBotInfoIV
	query := `
		SELECT
			*
		FROM
			line_works_bot_info_iv
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineWorksBotInfoIV, query, guildID)
	return &lineWorksBotInfoIV, err
}

func (r *Repository) UpdateLineWorksBotInfoIV(ctx context.Context, lineWorksBotInfoIV *LineWorksBotInfoIV) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineWorksBotInfoIV.LineWorksClientIDIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_client_id_iv = :line_works_client_id_iv")
	}
	if len(lineWorksBotInfoIV.LineWorksClientSecretIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_client_secret_iv = :line_works_client_secret_iv")
	}
	if len(lineWorksBotInfoIV.LineWorksServiceAccountIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_service_account_iv = :line_works_service_account_iv")
	}
	if len(lineWorksBotInfoIV.LineWorksPrivateKeyIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_private_key_iv = :line_works_private_key_iv")
	}
	if len(lineWorksBotInfoIV.LineWorksDomainIDIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_domain_id_iv = :line_works_domain_id_iv")
	}
	if len(lineWorksBotInfoIV.LineWorksAdminIDIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_admin_id_iv = :line_works_admin_id_iv")
	}
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE
			line_works_bot_info_iv
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)

	_, err := r.db.NamedExecContext(ctx, query, lineWorksBotInfoIV)
	return err
}
