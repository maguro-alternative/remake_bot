package internal

import (
	"time"

	"github.com/lib/pq"
)

type LineWorksTokenJson struct {
	GuildID                       string `json:"guildId" db:"guild_id"`
	LineWorksClientID             string `json:"lineWorksClientID,omitempty" db:"line_works_client_id"`
	LineWorksClientSecret         string `json:"lineWorksClientSecret,omitempty" db:"line_works_client_secret"`
	LineWorksServiceAccount       string `json:"lineWorksServiceAccount,omitempty" db:"line_works_service_account"`
	LineWorksPrivateKey           string `json:"lineWorksPrivateKey,omitempty" db:"line_works_private_key"`
	LineWorksDomainID             string `json:"lineWorksDomainID,omitempty" db:"line_works_domain_id"`
	LineWorksAdminID              string `json:"lineWorksAdminID,omitempty" db:"line_works_admin_id"`
	LineWorksBotID                string `json:"lineWorksBotID,omitempty" db:"line_works_bot_id"`
	LineWorksBotSecret            string `json:"lineWorksBotSecret,omitempty" db:"line_works_bot_secret"`
	LineWorksGroupID              string `json:"lineWorksGroupID,omitempty" db:"line_works_group_id"`
	DefaultChannelID              string `json:"defaultChannelId,omitempty" db:"default_channel_id"`
	DebugMode                     bool   `json:"debugMode,omitempty" db:"debug_mode"`
	LineWorksClientIDDelete       bool   `json:"lineWorksClientIdDelete,omitempty"`
	LineWorksClientSecretDelete   bool   `json:"lineWorksClientSecretDelete,omitempty"`
	LineWorksServiceAccountDelete bool   `json:"lineWorksServiceAccountDelete,omitempty"`
	LineWorksPrivateKeyDelete     bool   `json:"lineWorksPrivateKeyDelete,omitempty"`
	LineWorksDomainIDDelete       bool   `json:"lineWorksDomainIdDelete,omitempty"`
	LineWorksAdminIDDelete        bool   `json:"lineWorksAdminIdDelete,omitempty"`
	LineWorksBotIDDelete          bool   `json:"lineWorksBotIdDelete,omitempty"`
	LineWorksBotSecretDelete      bool   `json:"lineWorksBotSecretDelete,omitempty"`
	LineWorksGroupIDDelete        bool   `json:"lineWorksGroupIdDelete,omitempty"`
}

type LineWorksBot struct {
	GuildID                       string        `db:"guild_id"`
	LineWorksBotToken			 pq.ByteaArray `db:"line_works_bot_token"`
	LineWorksRefreshToken		 pq.ByteaArray `db:"line_works_refresh_token"`
	LineWorksGroupID              pq.ByteaArray `db:"line_works_group_id"`
	LineWorksBotID                pq.ByteaArray `db:"line_works_bot_id"`
	LineWorksBotSecret            pq.ByteaArray `db:"line_works_bot_secret"`
	RefreshTokenExpiresAt		 time.Time     `db:"refresh_token_expires_at"`
	DefaultChannelID              string        `db:"default_channel_id"`
	DebugMode					 bool          `db:"debug_mode"`
}

type LineWorksBotIv struct {
	GuildID                       string        `db:"guild_id"`
	LineWorksBotTokenIv           pq.ByteaArray `db:"line_works_bot_token_iv"`
	LineWorksRefreshTokenIv       pq.ByteaArray `db:"line_works_refresh_token_iv"`
	LineWorksGroupIDIv            pq.ByteaArray `db:"line_works_group_id_iv"`
	LineWorksBotIDIv              pq.ByteaArray `db:"line_works_bot_id_iv"`
	LineWorksBotSecretIv          pq.ByteaArray `db:"line_works_bot_secret_iv"`
}

type LineWorksBotInfo struct {
	GuildID                       string        `db:"guild_id"`
	LineWorksClientID             pq.ByteaArray `db:"line_works_client_id"`
	LineWorksClientSecret         pq.ByteaArray `db:"line_works_client_secret"`
	LineWorksServiceAccount       pq.ByteaArray `db:"line_works_service_account"`
	LineWorksPrivateKey           pq.ByteaArray `db:"line_works_private_key"`
	LineWorksDomainID             pq.ByteaArray `db:"line_works_domain_id"`
	LineWorksAdminID              pq.ByteaArray `db:"line_works_admin_id"`
}

type LineWorksBotInfoIv struct {
	GuildID                       string        `db:"guild_id"`
	LineWorksClientIDIv           pq.ByteaArray `db:"line_works_client_id_iv"`
	LineWorksClientSecretIv       pq.ByteaArray `db:"line_works_client_secret_iv"`
	LineWorksServiceAccountIv     pq.ByteaArray `db:"line_works_service_account_iv"`
	LineWorksPrivateKeyIv         pq.ByteaArray `db:"line_works_private_key_iv"`
	LineWorksDomainIDIv           pq.ByteaArray `db:"line_works_domain_id_iv"`
	LineWorksAdminIDIv            pq.ByteaArray `db:"line_works_admin_id_iv"`
}
