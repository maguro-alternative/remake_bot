/*
サーバーの権限設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT): 権限の種類 (line, line_bot, vc, webhook)
    permission_code (TEXT): DIscord上での権限コード
    permission (TEXT): 権限レベル
*/
CREATE TABLE IF NOT EXISTS permissions_code (
    guild_id TEXT NOT NULL,
    type TEXT NOT NULL,
    permission_code TEXT NOT NULL,
    permission TEXT NOT NULL,
    PRIMARY KEY(guild_id)
);

/*
サーバーの権限設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT): 権限の種類 (line, line_bot, vc, webhook)
    target_type (TEXT): 対象の種類 (user, role)
    target_id (TEXT): 対象ID (ユーザーID、ロールID)
    permission (TEXT): 権限レベル
*/
CREATE TABLE IF NOT EXISTS permissions_id (
    guild_id TEXT NOT NULL,
    type TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_id TEXT NOT NULL,
    permission TEXT NOT NULL,
    PRIMARY KEY(guild_id)
);

/*
DiscordからLINEへのメッセージ送信設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    channel_id (TEXT): チャンネルID
    ng (BOOLEAN): 送信NGのチャンネルか
    bot_message (BOOLEAN): Botのメッセージを送信するか
*/
CREATE TABLE IF NOT EXISTS line_channels (
    channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    ng BOOLEAN NOT NULL,
    bot_message BOOLEAN NOT NULL,
    PRIMARY KEY(channel_id)
);

/*
LINEへ送信しないメッセージの種類を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT PRIMARY KEY): メッセージの種類
*/
CREATE TABLE IF NOT EXISTS line_ng_types (
    guild_id TEXT NOT NULL,
    type TEXT NOT NULL,
    PRIMARY KEY(guild_id, type)
);

/*
LINEへ送信しないDiscordユーザーを保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    user_id (TEXT PRIMARY KEY): ユーザーID
*/
CREATE TABLE IF NOT EXISTS line_ng_discord_users (
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY(guild_id, user_id)
);

/*
LINEへ送信しないDiscordロールを保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    role_id (TEXT PRIMARY KEY): ロールID
*/
CREATE TABLE IF NOT EXISTS line_ng_discord_roles (
    guild_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    PRIMARY KEY(guild_id, role_id)
);

CREATE TABLE IF NOT EXISTS vc_signal (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    send_signal BOOLEAN NOT NULL,
    send_channel_id TEXT NOT NULL,
    join_bot BOOLEAN NOT NULL,
    everyone_mention BOOLEAN NOT NULL,
    PRIMARY KEY(vc_channel_id)
);

CREATE TABLE IF NOT EXISTS vc_signal_ng (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, user_id)
);

CREATE TABLE IF NOT EXISTS vc_signal_ng_role (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, role_id)
);

CREATE TABLE IF NOT EXISTS vc_signal_mention_user (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, user_id)
);

CREATE TABLE IF NOT EXISTS vc_signal_mention_role (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, role_id)
);

CREATE TABLE IF NOT EXISTS webhook (
    id SERIAL PRIMARY KEY,
    guild_id TEXT NOT NULL,
    webhook_id TEXT NOT NULL,
    subscription_type TEXT NOT NULL,
    subscription_id TEXT NOT NULL,
    last_posted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS webhook_mention_users (
    id INTEGER PRIMARY KEY,
    user_id TEXT NOT NULL,
    PRIMARY KEY(id, user_id),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS webhook_mention_roles (
    id INTEGER PRIMARY KEY,
    role_id TEXT NOT NULL,
    PRIMARY KEY(id, role_id),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS webhook_ng_or_words (
    id INTEGER PRIMARY KEY,
    word TEXT NOT NULL,
    PRIMARY KEY(id, word),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS webhook_ng_and_words (
    id INTEGER PRIMARY KEY,
    word TEXT NOT NULL,
    PRIMARY KEY(id, word),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS webhook_search_or_words (
    id INTEGER PRIMARY KEY,
    word TEXT NOT NULL,
    PRIMARY KEY(id, word),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS webhook_search_and_words (
    id INTEGER PRIMARY KEY,
    word TEXT NOT NULL,
    PRIMARY KEY(id, word),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS webhook_mention_or_words (
    id INTEGER PRIMARY KEY,
    word TEXT NOT NULL,
    PRIMARY KEY(id, word),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS webhook_mention_and_words (
    id INTEGER PRIMARY KEY,
    word TEXT NOT NULL,
    PRIMARY KEY(id, word),
    FOREIGN KEY(id) REFERENCES webhook(id)
);

CREATE TABLE IF NOT EXISTS line_bot (
    guild_id TEXT NOT NULL,
    line_notify_token BYTEA,
    line_bot_token BYTEA,
    line_bot_secret BYTEA,
    line_group_id BYTEA,
    line_client_id BYTEA,
    line_client_secret BYTEA,
    default_channel_id TEXT,
    debug_mode BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY(guild_id)
);
