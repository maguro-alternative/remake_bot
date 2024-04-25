/*
サーバーの権限設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT): 権限の種類 (line_post_discord_channel, line_bot, vc, webhook)
    code (BIGINT): Discord上での権限コード
*/
CREATE TABLE IF NOT EXISTS permissions_code (
    guild_id TEXT NOT NULL,
    type TEXT NOT NULL,
    code BIGINT NOT NULL,
    PRIMARY KEY(guild_id, type)
);

/*
サーバーのユーザー権限設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT): 権限の種類 (line_post_discord_channel, line_bot, vc, webhook)
    user_id (TEXT): 対象ID (ユーザーID)
    permission (TEXT): 権限レベル(read, write, admin)
*/
CREATE TABLE IF NOT EXISTS permissions_user_id (
    guild_id TEXT NOT NULL,
    type TEXT NOT NULL,
    user_id TEXT NOT NULL,
    permission TEXT NOT NULL,
    PRIMARY KEY(guild_id, type, user_id)
);

/*
サーバーのロール権限設定を保存するテーブル

カラム:

    guild_id (TEXT PRIMARY KEY): サーバーID
    type (TEXT): 権限の種類 (line_post_discord_channel, line_bot, vc, webhook)
    user_id (TEXT): 対象ID (ロールID)
    permission (TEXT): 権限レベル(read, write, admin)
*/
CREATE TABLE IF NOT EXISTS permissions_role_id (
    guild_id TEXT NOT NULL,
    type TEXT NOT NULL,
    role_id TEXT NOT NULL,
    permission TEXT NOT NULL,
    PRIMARY KEY(guild_id, type, role_id)
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
LINEへ送信しないDiscordユーザーを保存するテーブル

カラム:

    channel_id (TEXT): チャンネルID
    guild_id (TEXT PRIMARY KEY): サーバーID
    id (TEXT PRIMARY KEY): ID
*/
CREATE TABLE IF NOT EXISTS line_ng_discord_user_id (
    channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY(channel_id, user_id)
);

/*
LINEへ送信しないDiscordロールを保存するテーブル

カラム:

    channel_id (TEXT): チャンネルID
    guild_id (TEXT PRIMARY KEY): サーバーID
    id (TEXT PRIMARY KEY): ID
*/
CREATE TABLE IF NOT EXISTS line_ng_discord_role_id (
    channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    PRIMARY KEY(channel_id, role_id)
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
指定されたユーザーがボイスチャンネルに参加した場合通知しない

カラム:

    vc_channel_id (TEXT): ボイスチャンネルID
    guild_id (TEXT): サーバーID
    id (TEXT): ID
*/

CREATE TABLE IF NOT EXISTS vc_signal_ng_user_id (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, user_id)
);

/*
指定されたロールがボイスチャンネルに参加した場合通知しない

カラム:

    vc_channel_id (TEXT): ボイスチャンネルID
    guild_id (TEXT): サーバーID
    id (TEXT): ID
*/


CREATE TABLE IF NOT EXISTS vc_signal_ng_role_id (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, role_id)
);

/*
ボイスチャンネルの通知の際にメンションするユーザー、ロールを保存するテーブル

カラム:

    vc_channel_id (TEXT): ボイスチャンネルID
    guild_id (TEXT): サーバーID
    user_id (TEXT): ユーザーID
*/

CREATE TABLE IF NOT EXISTS vc_signal_mention_user_id (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, user_id)
);

CREATE TABLE IF NOT EXISTS vc_signal_mention_role_id (
    vc_channel_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    PRIMARY KEY(vc_channel_id, role_id)
);

CREATE TABLE IF NOT EXISTS webhook (
    webhook_serial_id SERIAL,
    guild_id TEXT NOT NULL,
    webhook_id TEXT NOT NULL,
    subscription_type TEXT NOT NULL,
    subscription_id TEXT NOT NULL,
    last_posted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(webhook_serial_id)
);

CREATE TABLE IF NOT EXISTS webhook_user_mention (
    webhook_serial_id INTEGER,
    user_id TEXT NOT NULL,
    PRIMARY KEY(webhook_serial_id, user_id),
    FOREIGN KEY(webhook_serial_id) REFERENCES webhook(webhook_serial_id)
);

CREATE TABLE IF NOT EXISTS webhook_role_mention (
    webhook_serial_id INTEGER,
    role_id TEXT NOT NULL,
    PRIMARY KEY(webhook_serial_id, role_id),
    FOREIGN KEY(webhook_serial_id) REFERENCES webhook(webhook_serial_id)
);

/*ng_or ng_and search_or search_and mention_or mention_and*/
CREATE TABLE IF NOT EXISTS webhook_word (
    webhook_serial_id INTEGER,
    conditions TEXT NOT NULL,
    word TEXT NOT NULL,
    PRIMARY KEY(webhook_serial_id, word),
    FOREIGN KEY(webhook_serial_id) REFERENCES webhook(webhook_serial_id)
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

/*DROP TABLE IF EXISTS permission_code, permission_id, line_post_discord_channel, line_ng_discord_message_type, line_ng_discord_user_id, line_ng_discord_role_id, vc_signal_channel, vc_signal_ng_id, vc_signal_mention_id, webhook, webhook_mention, webhook_word, line_bot, line_bot_iv CASCADE;*/
