# discordfastリプレイス

https://github.com/maguro-alternative/discordfast

上記のリポジトリのリプレイスです。

# 進捗

- [x] LINE→Discordのメッセージの送信
- [x] Discord→LINEのメッセージの送信
- [x] LINE→Discordの画像の送信
- [x] Discord→LINEの画像の送信
- [x] LINE→Discordの動画の送信
- [x] Discord→LINEの動画の送信
- [x] LINE→Discordの音声の送信
- [x] Discord→LINEの音声の送信
- [x] LINE→Discordのスタンプの送信
- [x] Discord→LINEのスタンプの送信
- [] LINE→Discordのファイルの送信
- [] Discord→LINEのファイルの送信

# ディレクトリ構造

<details>
    <summary>ディレクトリ概要(折り畳み)</summary>

```plaintext
.
├── bot                         // DiscordBotを動かすためのディレクトリ
│   ├── cogs                    // DiscordBotのコグ
│   │   └── on_message_create   // メッセージが送信されたときのデータベースの操作
|   ├── commands                // スラッシュコマンド
│   ├── config                  // 環境変数設定ファイル
│   └── main.go
├── core                        // サーバーとBotを動かすためのディレクトリ
│   ├── config                  // 環境変数設定ファイル
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
|   ├── main.go
│   └── schema.sql              // データベースのスキーマ
├── fixtures                    // データベースのテスト用のフィクスチャ
├── pkg                         // 共通のパッケージ
│   ├── crypto                  // 暗号化関連のパッケージ
│   ├── db                      // データベース関連のパッケージ
│   ├── line                    // LINEBot関連のパッケージ
│   └── youtube                 // YouTube関連のパッケージ
├── web                         // Webサーバーを動かすためのディレクトリ
│   ├── components              // Webサーバーのコンポーネント
│   ├── config                  // 環境変数設定ファイル
│   ├── handler                 // Webサーバーのハンドラ
│   ├── middleware              // Webサーバーのミドルウェア
│   ├── service                 // Webサーバーのサービス
│   ├── shared                  // Webサーバー内での共通のパッケージ
│   └── templates               // WebサーバーのHTMLテンプレート
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

</details>

<details>
    <summary>全体(長すぎるので折り畳み)</summary>

```plaintext
.
├── bot
│   ├── cogs
│   │   ├── on_message_create                       // メッセージが送信されたときのデータベースの操作
│   │   │   ├── entity.go
│   │   │   ├── repository_test.go
│   │   │   └── repository.go
│   │   ├── cog_handler.go                          // ここでcogを登録
│   │   ├── on_message_create.go                    // discord内でメッセージが送信されたときの処理
│   │   └── vc_signal.go                            // ボイスチャンネルのステータス変化時の処理
|   ├── commands
|   |   ├── command_handler.go                      // ここでコマンドを登録
|   |   └── ping.go                                 // pingコマンド
│   ├── config                                      // 環境変数設定ファイル
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   └── main.go                                     // Botメイン関数
├── core
│   ├── config                                      // 環境変数設定ファイル
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
|   ├── main.go
│   └── schema.sql                                  // データベースのスキーマ
├── fixtures
│   ├── fixtures.go
│   ├── line_bot_iv.go
│   ├── line_bot.go
|   ├── line_ng_discord_id.go
|   ├── line_ng_discord_message_type.go
│   ├── line_post_discord_chennel.go
│   ├── permissions_code.go
│   ├── permissions_id.go
│   ├── vc_signal_channel.go
│   ├── vc_signal_mention_id.go
│   ├── vc_signal_ng_id.go
│   ├── webhook_mention.go
│   ├── webhook_word.go
│   └── webhook.go
├── pkg
│   ├── crypto
│   │   ├── aes.go
│   │   └── aes_test.go
│   ├── db
│   │   ├── db.go
│   │   └── db_test.go
│   ├── line
│   │   ├── get_bot_info.go
│   │   ├── get_friend_count.go
│   │   ├── get_group_count.go
│   │   ├── get_message_content.go
│   │   ├── get_profile.go
│   │   ├── get_pushlimit.go
│   │   ├── get_totalpush_count.go
│   │   ├── line_message.go
│   │   ├── line_notify.go
│   │   └── line.go
│   └── youtube
│       ├── create_client_secret.go
│       ├── create_oauth2.go
│       └── youtube.go
├── web
│   ├── components
│   │   ├── discord_account_pop.go
│   │   ├── entity.go
│   │   ├── line_account_pop.go
│   │   ├── line_post_discord_chennel.go
│   │   ├── linetoken.go
│   │   └── submittag.go
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   ├── handler
│   │   ├── api
│   │   │   ├── group
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   └── group.go
│   │   |   ├── line_post_discord_chennel
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   └── line_post_discord_chennel.go
│   │   │   ├── line_bot
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── hmac.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   ├── entity.go
│   │   │   │   └── line_bot.go
│   │   │   └── linetoken
│   │   │       ├── internal
│   │   │       │   ├── entity.go
│   │   │       │   ├── repository_test.go
│   │   │       │   └── repository.go
│   │   │       └── linetoken.go
│   │   ├── callback
│   │   │   ├── discord_callback
│   │   │   │   └── callback.go
│   │   │   └── line_callback
│   │   │       └── callback.go
│   │   ├── login
│   │   │   ├── discord_login
│   │   │   │   └── discord_login.go
│   │   │   └── line_login
│   │   │       ├── internal
│   │   │       │   ├── entity.go
│   │   │       │   ├── repository_test.go
│   │   │       │   └── repository.go
│   │   │       └── line_login.go
│   │   ├── logout
│   │   │   ├── discord_logout
│   │   │   │   └── discord_logout.go
│   │   │   └── line_logout
│   │   │       └── line_logout.go
│   │   ├── views
│   │   │   ├── group
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   └── group.go
│   │   │   ├── guildid
│   │   │   │   ├── line_post_discord_chennel
│   │   │   │   │   ├── internal
│   │   │   │   │   │   ├── entity.go
│   │   │   │   │   │   ├── repository_test.go
│   │   │   │   │   │   └── repository.go
│   │   │   │   │   └── line_post_discord_chennel.go
│   │   │   │   ├── line_token
│   │   │   │   │   ├── internal
│   │   │   │   │   │   ├── entity.go
│   │   │   │   │   │   ├── repository_test.go
│   │   │   │   │   │   └── repository.go
│   │   │   │   │   └── linetoken.go
│   │   │   │   └── guildid.go
│   │   │   ├── guilds
│   │   │   │   └── guilds.go
│   │   │   └── index.go
│   ├── middleware
│   │   └── middleware.go
│   ├── service
│   │   ├── discord_oauth2.go
│   │   └── index.go
│   ├── shared
│   │   ├── permission
│   │   │   ├── internal
│   │   │   │   ├── entity.go
│   │   │   │   ├── repository_test.go
│   │   │   │   └── repository.go
│   │   │   ├── check_discord_permission.go
│   │   │   └── check_line_permission.go
│   │   └── session
│   │       ├── getoauth
│   │       │   ├── get_discord_oauth.go
│   │       │   └── get_line_oauth.go
│   │       └── model
│   │           └── entity.go
│   └── templates
│       ├── static
│       │   ├── img
│       │   │   └── logo.png
│       │   └── js
│       │       ├── group.js
│       │       ├── line_post_discord_chennel.js
│       │       ├── linetoken.js
│       │       └── popover.js
│       ├── views
│       │   ├── group
│       │   │   └── group.html
│       │   ├── guildid
│       │   │   ├── line_post_discord_chennel.html
│       │   │   └── linetoken.html
│       │   ├── guilds
│       │   │   └── guilds.html
|       │   ├── login
│       │   │   └── line_login.html
│       │   └── guildid.html
│       ├── index.html
│       └── layout.html
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

</details>

# データベース

太文字は主キー  
複数のカラムが太文字になっている場合は複合主キー  

<details>
    <summary>permissions_code</summary>

サーバーの権限設定を保存するテーブル  
権限コードをすべて満たすユーザーが設定変更を行える  

|カラム名|型|説明|
|---|---|---|
|**guild_id**|TEXT|DiscordサーバーのID|
|**type**|TEXT|権限の種類 (line_post_discord_channel, line_bot, vc_signal, webhook)|
|code|BIGINT|Discordの権限コード、詳細は[こちら](https://discord.com/developers/docs/topics/permissions)|

</details>

<details>
    <summary>permissions_id</summary>

サーバーの権限設定を保存するテーブル
ここに保存されているユーザーやロールは、```permission```と同じ権限を持っているということになる

|カラム名|型|説明|
|---|---|---|
|**guild_id**|TEXT|DiscordサーバーのID|
|**type**|TEXT|権限の種類 (line_post_discord_channel, line_bot, vc_signal, webhook)|
|**target_type**|TEXT|権限の対象の種類 (role, user)|
|**target_id**|TEXT|権限の対象ID (ユーザーID、ロールID)|
|permission|TEXT|権限レベル(read, write, all)|

</details>

<details>
    <summary>line_post_discord_channel</summary>

DiscordからLINEグループにメッセージを送信するための設定を保存するテーブル

|カラム名|型|説明|
|---|---|---|
|**channel_id**|TEXT|DiscordのチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|ng|BOOLEAN|LINEに送信NGのチャンネルかどうか|
|bot_message|BOOLEAN|DiscordBotのメッセージを送信するかどうか|

</details>

<details>
    <summary>line_ng_discord_message_type</summary>

LINEに送信NGのDiscordメッセージの種類を保存するテーブル  
discordgo.MessageTypeで使用されている定数(0~23)と同じ値を使用する  

|カラム名|型|説明|
|---|---|---|
|**channel**|TEXT|DiscordのチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|**type**|INTEGER|メッセージの種類(ピン止め、スレッド、スレッドの返信)|

</details>

<details>
    <summary>line_ng_discord_id</summary>

LINEへ送信しないDiscordユーザー、ロールを保存するテーブル  
ここに保存されているユーザー、ロールを持つユーザーはLINEにメッセージが送信されない  

|カラム名|型|説明|
|---|---|---|
|**channel**|TEXT|DiscordのチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|**id**|TEXT|ID|
|id_type|TEXT|IDの種類 (user, role)|

</details>

<details>
    <summary>vc_signal_channel</summary>

ボイスチャンネル入退出の通知設定を保存するテーブル

|カラム名|型|説明|
|---|---|---|
|**vc_channel_id**|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|send_signal|BOOLEAN|L通知を送信するかどうか|
|send_channel_id|TEXT|通知を送信するチャンネルID|
|join_bot|BOOLEAN|ボイスチャンネルにBotが入室したときの通知を送信するかどうか|
|everyone_mention|BOOLEAN|通知を送信するときに@everyoneを使用するかどうか|

</details>

<details>
    <summary>vc_signal_ng_id</summary>

指定されたユーザー、ロールがボイスチャンネルに参加した場合通知しない

|カラム名|型|説明|
|---|---|---|
|**vc_channel_id**|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|**id**|TEXT|ID|
|id_type|TEXT|IDの種類 (user, role)|

</details>

<details>
    <summary>vc_signal_mention_id</summary>

ボイスチャンネルの通知の際にメンションするユーザー、ロールを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|**vc_channel_id**|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|**id**|TEXT|ID|
|id_type|TEXT|IDの種類 (user, role)|

</details>

<details>
    <summary>webhook</summary>

DiscordのWebhookの送信設定を保存するテーブル

|カラム名|型|説明|
|---|---|---|
|**webhook_serial_id**|SERIAL|Webhookの投稿内容の識別ID|
|guild_id|TEXT|DiscordのサーバーID|
|webhook_id|TEXT|WebhookのID|
|subscription_type|TEXT|読み取るもの(YouTube,NicoNico)|
|subscription_id|TEXT|上記のサービスで投稿者を識別できるもの|
|last_posted_at|TIMESTAMP|最後に投稿した日時|

</details>

<details>
    <summary>webhook_mention</summary>

Webhookの送信時にメンションするユーザー、ロールを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|**webhook_serial_id**|SERIAL|Webhookの投稿内容の識別ID|
|**id**|TEXT|ID|
|id_type|TEXT|IDの種類 (user, role)|

</details>

<details>
    <summary>webhook_word</summary>

Webhookの送信時に特定の単語が含まれていた場合にメンションするユーザー、ロールを保存するテーブル  
Twitter運用時使用していたが現在死に要素  
conditionsは投稿時の条件を示す(ng_orはいずれかの単語が含まれていれば投稿しない。search_andは全ての単語が含まれていれば投稿する。mention_orはいずれかの単語が含まれていればメンションする。)

|カラム名|型|説明|
|---|---|---|
|**webhook_serial_id**|SERIAL|Webhookの投稿内容の識別ID|
|**word**|TEXT|メンションする単語|
|conditions|TEXT|投稿時の条件(ng_or ng_and search_or search_and mention_or mention_and)|

</details>

<details>
    <summary>line_bot</summary>

LINEBotの設定を保存するテーブル  
LINEBotのアクセストークン、チャンネルシークレットなどをAES暗号化して保存する

|カラム名|型|説明|
|---|---|---|
|**guild_id**|TEXT|DiscordのサーバーID|
|line_notify_token|BYTEA|LINE Notifyのトークン(AESで暗号化)|
|line_bot_token|BYTEA|LINEBotのアクセストークン(AESで暗号化)|
|line_bot_secret|BYTEA|LINEBotのチャンネルシークレット(AESで暗号化)|
|line_group_id|BYTEA|LINEのグループID(AESで暗号化)|
|line_client_id|BYTEA|LINEのクライアントID(AESで暗号化)|
|line_client_secret|BYTEA|LINEのクライアントシークレット(AESで暗号化)|
|default_channel_id|TEXT|LINEに送信するチャンネルID|
|debug_mode|BOOLEAN|デバッグモードかどうか(オンにするとLINEグループにメッセージを送信するたびLINEのグループIDが返ってくる)|

</details>

<details>
    <summaty>line_bot_iv</summary>

LINEBotの復号化に使用するIVを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|**guild_id**|TEXT|DiscordのサーバーID|
|line_notify_token_iv|BYTEA|LINE NotifyトークンのIV|
|line_bot_token_iv|BYTEA|LINEBotのアクセストークンのIV|
|line_bot_secret_iv|BYTEA|LINEBotのチャンネルシークレットのIV|
|line_group_id_iv|BYTEA|LINEのグループIDのIV|
|line_client_id_iv|BYTEA|LINEのクライアントIDのIV|
|line_client_secret_iv|BYTEA|LINEのクライアントシークレットのIV|

</details>
