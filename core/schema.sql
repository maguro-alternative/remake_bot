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
    guild_id TEXT NOT NULL,
    channel_id TEXT NOT NULL,
    ng BOOLEAN NOT NULL,
    bot_message BOOLEAN NOT NULL,
    PRIMARY KEY(guild_id)
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
