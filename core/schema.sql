/*
サーバーの権限設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT): 権限の種類 (line, line_bot, vc, webhook)
    code (BIGINT): Discord上での権限コード
*/
CREATE TABLE IF NOT EXISTS permissions_code (
    guild_id TEXT NOT NULL,
    type TEXT NOT NULL,
    code BIGINT NOT NULL,
    PRIMARY KEY(guild_id, type)
);

/*
サーバーの権限設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT): 権限の種類 (line, line_bot, vc, webhook)
    target_type (TEXT): 対象の種類 (user, role)
    target_id (TEXT): 対象ID (ユーザーID、ロールID)
    permission (TEXT): 権限レベル(read, write, admin)
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
CREATE TABLE IF NOT EXISTS line_post_discord_channel (
    channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    ng BOOLEAN NOT NULL,
    bot_message BOOLEAN NOT NULL,
    PRIMARY KEY(channel_id)
);

/*
LINEへ送信しないメッセージの種類を保存するテーブル

カラム:

    channel_id (TEXT PRIMARY KEY): チャンネルID
    guild_id (TEXT PRIMARY KEY): サーバーID
    type (INTEGER PRIMARY KEY): メッセージの種類(ピン止め、スレッド、スレッドの返信)
*/
CREATE TABLE IF NOT EXISTS line_ng_discord_message_type (
    channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    type INTEGER NOT NULL,
    PRIMARY KEY(channel_id, type)
);

/*
LINEへ送信しないDiscordユーザー、ロールを保存するテーブル

カラム:

    channel_id (TEXT): チャンネルID
    guild_id (TEXT PRIMARY KEY): サーバーID
    id (TEXT PRIMARY KEY): ID
    id_type (TEXT PRIMARY KEY): ユーザーIDの種類 (user, role)
*/
CREATE TABLE IF NOT EXISTS line_ng_discord_id (
    channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    id TEXT NOT NULL,
    id_type TEXT NOT NULL,
    PRIMARY KEY(channel_id, id)
);

/*
ボイスチャンネルの通知設定を保存するテーブル

カラム:

    vc_channel_id (TEXT PRIMARY KEY): ボイスチャンネルID
    guild_id (TEXT): サーバーID
    send_signal (BOOLEAN): 通知を送信するか
    send_channel_id (TEXT): 通知を送信するチャンネルID
    join_bot (BOOLEAN): Botの参加を通知するか
    everyone_mention (BOOLEAN): @everyoneを通知するか
*/

CREATE TABLE IF NOT EXISTS vc_signal_channel (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    send_signal BOOLEAN NOT NULL,
    send_channel_id TEXT NOT NULL,
    join_bot BOOLEAN NOT NULL,
    everyone_mention BOOLEAN NOT NULL,
    PRIMARY KEY(vc_channel_id)
);

/*
指定されたユーザー、ロールがボイスチャンネルに参加した場合通知しない

カラム:

    vc_channel_id (TEXT): ボイスチャンネルID
    guild_id (TEXT): サーバーID
    id_type (TEXT PRIMARY KEY): ユーザーIDの種類 (user, role)
    id (TEXT): ID
*/

CREATE TABLE IF NOT EXISTS vc_signal_ng_id (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    id_type TEXT NOT NULL,
    id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, id)
);

/*
ボイスチャンネルの通知の際にメンションするユーザー、ロールを保存するテーブル

カラム:

    vc_channel_id (TEXT): ボイスチャンネルID
    guild_id (TEXT): サーバーID
    id_type (TEXT PRIMARY KEY): ユーザーIDの種類 (user, role)
    id (TEXT): ユーザーID
*/

CREATE TABLE IF NOT EXISTS vc_signal_mention_id (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    id_type TEXT NOT NULL,
    id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, id)
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

CREATE TABLE IF NOT EXISTS webhook_mention (
    webhook_serial_id INTEGER PRIMARY KEY,
    id_type TEXT NOT NULL,
    id TEXT NOT NULL,
    PRIMARY KEY(webhook_serial_id, id),
    FOREIGN KEY(webhook_serial_id) REFERENCES webhook(id)
);

/*ng_or ng_and search_or search_and mention_or mention_and*/
CREATE TABLE IF NOT EXISTS webhook_word (
    id INTEGER PRIMARY KEY,
    conditions TEXT NOT NULL,
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

CREATE TABLE IF NOT EXISTS line_bot_iv (
    guild_id TEXT NOT NULL,
    line_notify_token_iv BYTEA,
    line_bot_token_iv BYTEA,
    line_bot_secret_iv BYTEA,
    line_group_id_iv BYTEA,
    line_client_id_iv BYTEA,
    line_client_secret_iv BYTEA,
    PRIMARY KEY(guild_id)
);

/*DROP TABLE IF EXISTS permission_code, permission_id, line_post_discord_channel, line_ng_discord_message_type, line_ng_discord_id, vc_signal_channel, vc_signal_ng_id, vc_signal_mention_id, webhook, webhook_mention, webhook_word, line_bot, line_bot_iv CASCADE;*/
