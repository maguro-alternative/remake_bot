# discordfastリプレイス

https://github.com/maguro-alternative/discordfast

上記のリポジトリのリプレイスです。

# 機能
## ボイスチャンネル入退室通知

<details>
    <summary>画像(折り畳み)</summary>

- 入退出  
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/0b4b012d-0b24-4803-9007-6ce2e652fa02)
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/4c1f353b-b598-4684-a7e8-fe1b64c6501d)
- 画面共有、カメラ配信  
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/b10167e6-1633-4491-95e8-48750c4c7299)

</details>

## LINE連携
利用の際は下記の画像にもあるように、WebUI上のフォームにトークンなどを入力する必要があります。
<details>
    <summary>画像(折り畳み)</summary>

- メッセージ  
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/814b64e4-4007-42dd-aad8-0134e80801a2)
- 動画、スタンプ  
![image](https://github.com/maguro-alternative/discordfast/raw/master/img/134639.png)

</details>

## VoiceVox読み上げ
<details>
    <summary>画像(折り畳み)</summary>

- テキスト読み上げ  
![image](https://github.com/maguro-alternative/discordfast/raw/master/img/zundacommand.png)

</details>

## Webhook
<details>
    <summary>画像(折り畳み)</summary>

- YouTubeの新着動画  
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/de60175c-31b8-4c79-ad8c-952a2bd62191)

</details>


# 使い方
以下の環境変数を設定してください。  
PostgresSQLのデータベースを使用します。  
Windowsの場合```go.mod```のディレクトリに```.env```ファイルを作成し、記述してください。  

```plaintext
# DiscordのBotのトークン
DISCORD_BOT_TOKEN=your_discord_bot_token
DISCORD_CLIENT_ID=your_discord_client_id
DISCORD_CLIENT_SECRET=your_discord_client_secret
DISCORD_SCOPE=identify%20guilds.join%20guilds%20guilds.members.read%20applications.builds.read%20connections

# DB_は現時点のコードでは使用されていません(テストコードでのみ運用)
DB_NAME=your_db_name
DB_HOST=your_db_host
DB_PASSWORD=your_db_password
DB_PORT=your_db_port
DB_USER=your_db_user
# データベースのURL
DATABASE_URL=

PORT=5000
# AES暗号化のためのキー(128bit,00~FFの16進数32文字)
PRIVATE_KEY=your_private_key
# サーバーのURL(各OAuth2のコールバック用)
SERVER_URL=your_server_url

# gorrila/sessionsのためのキー
SESSION_SECRET=your_session_secret

# YouTubeのAPIキー
YOUTUBE_API_KEY=your_youtube_api_key
YOUTUBE_ACCESS_TOKEN=your_youtube_access_token
YOUTUBE_CLIENT_ID=your_youtube_client_id
YOUTUBE_CLIENT_SECRET=your_youtube_client_secret
YOUTUBE_PROJECT_ID=your_youtube_project_id
YOUTUBE_REFRESH_TOKEN=your_youtube_refresh_token
YOUTUBE_TOKEN_EXPIRY=your_youtube_token_expiry

# Web版VoiceVoxのAPIキー
VOICEVOX_KEY=your_voicevox_key
```

以下のコマンドで依存関係を解決します。
```bash
go mod download
```

以下のコマンドでサーバーとBotを起動します。  
```bash
go run core/main.go
```

![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/0ae11602-a95a-4bee-b324-0afd1d6ecf57)  
起動時に所属サーバーを読み取り、自動的にデータベースに必要なものをInsertしてくれます。  
テーブルがない場合も自動的に作成されます。

![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/f6e8d05c-83a6-42ca-b41f-95b5b62642d8)  
サーバーが立ち上がったらアクセスしてみましょう。  
上記のようなページが表示されます。  

アカウントのところをクリックするとログインページへ遷移できます。  
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/8fcfa72a-aa57-4005-b227-63429a3f67d2)  
認可されると右のような画面になります。  
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/f7c6079b-b13c-43dc-9840-b893acbd439e)
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/2e3c69bf-e097-49a8-a1ad-711107e5f446)

認可を基に設定情報にアクセスできます。  
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/9819323a-e123-4dcc-a77c-df96335d570b)  

![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/6ac0cbcf-608d-4ad3-9fd8-2a4013307432)  

![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/57ae58ec-92fc-4e63-aaec-4f7a0ec4077d)

LINE連携を使用する場合は以下のフォームにトークンを入力する必要があります。
![image](https://github.com/maguro-alternative/remake_bot/assets/71870614/4fef03cd-4a18-4691-a388-c33cf333221c)

# ディレクトリ構造

<details>
    <summary>ディレクトリ概要(折り畳み)</summary>

```plaintext
.
├── bot                         // DiscordBotを動かすためのディレクトリ
│   ├── cogs                    // DiscordBotのコグ
|   ├── commands                // スラッシュコマンド
│   ├── config                  // 環境変数設定ファイル
│   ├── ffmpeg                  // 動画、音声の変換
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
│   ├── sharedtime              // botがVCに入室した時間を共有するためのパッケージ
│   └── youtube                 // YouTube関連のパッケージ
├── repository                  // データベース操作のリポジトリ
├── tasks                       // 定期的に行うタスク(Webhookの送信など)
├── testutil                    // テスト用のユーティリティ
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
│   │   ├── internal
│   │   │   └── entity.go
│   │   ├── cog_handler.go
│   │   ├── on_message_create.go
│   │   ├── on_message_create_test.go
│   │   ├── on_voice_state_update.go
│   │   └── on_voice_state_update_test.go
|   ├── commands
|   |   ├── command_handler.go
|   |   ├── command_handler_test.go
|   |   ├── ping_test.go
|   |   ├── ping.go
|   |   ├── voicevox_test.go
|   |   └── voicevox.go
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   ├── ffmpeg
│   │   ├── ffmpeg_test.go
│   │   └── ffmpeg.go
│   └── main.go
├── core
│   ├── config                                      // 環境変数設定ファイル
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
|   ├── main.go
│   └── schema.sql                                  // データベースのスキーマ
├── fixtures
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
├── repository
│   ├── line_bot_iv_test.go
│   ├── line_bot_iv.go
│   ├── line_bot_test.go
│   ├── line_bot.go
│   ├── line_ng_discord_message_type_test.go
│   ├── line_ng_discord_message_type.go
│   ├── line_ng_discord_role_id_test.go
│   ├── line_ng_discord_role_id.go
│   ├── line_ng_discord_user_id_test.go
│   ├── line_ng_discord_user_id.go
│   ├── line_post_discord_channel_test.go
│   ├── line_post_discord_channel.go
│   ├── permissions_code_test.go
│   ├── permissions_code.go
│   ├── permissions_role_id_test.go
│   ├── permissions_role_id.go
│   ├── permissions_user_id_test.go
│   ├── permissions_user_id.go
│   ├── repository.go
│   ├── vc_signal_channel_test.go
│   ├── vc_signal_channel.go
│   ├── vc_signal_mention_role_id_test.go
│   ├── vc_signal_mention_role_id.go
│   ├── vc_signal_mention_user_id_test.go
│   ├── vc_signal_mention_user_id.go
│   ├── vc_signal_ng_role_id_test.go
│   ├── vc_signal_ng_role_id.go
│   ├── vc_signal_ng_user_id_test.go
│   ├── vc_signal_ng_user_id.go
│   ├── webhook_role_mention_test.go
│   ├── webhook_role_mention.go
│   ├── webhook_user_mention_test.go
│   ├── webhook_user_mention.go
│   ├── webhook_word_test.go
│   ├── webhook_word.go
│   ├── webhook_test.go
│   └── webhook.go
├── tasks
│   ├── internal
│   │   ├── youtube_test.go
│   │   └── youtube.go
│   └── main.go
├── testutil
│   ├── files
│   │  └── video.xml
│   ├── fixtures
│   │   ├── fixtures.go
│   │   ├── line_bot_iv.go
│   │   ├── line_bot.go
│   |   ├── line_ng_discord_message_type.go
│   |   ├── line_ng_discord_role_id.go
│   |   ├── line_ng_discord_user_id.go
│   │   ├── line_post_discord_chennel.go
│   │   ├── permissions_code.go
│   │   ├── permissions_role_id.go
│   │   ├── permissions_user_id.go
│   │   ├── vc_signal_channel.go
│   │   ├── vc_signal_mention_role_id.go
│   │   ├── vc_signal_mention_user_id.go
│   │   ├── vc_signal_ng_role_id.go
│   │   ├── vc_signal_ng_user_id.go
│   │   ├── webhook_role_mention.go
│   │   ├── webhook_user_mention.go
│   │   ├── webhook_word.go
│   │   └── webhook.go
│   └── mock
│       ├── client.go
│       ├── repository.go
│       └── session.go
├── web
│   ├── components
│   │   ├── channel_select.go
│   │   ├── discord_account_pop.go
│   │   ├── entity.go
│   │   ├── line_account_pop.go
│   │   └── submittag.go
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   ├── handler
│   │   ├── api
│   │   │   ├── group
│   │   │   │   ├── internal
│   │   │   │   │   └── entity.go
│   │   │   │   ├── group_test.go
│   │   │   │   └── group.go
│   │   |   ├── line_post_discord_chennel
│   │   │   │   ├── internal
│   │   │   │   │   └── entity.go
│   │   │   │   ├── line_post_discord_channel_test.go
│   │   │   │   └── line_post_discord_chennel.go
│   │   │   ├── linebot
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── hmac.go
│   │   │   │   │   └── hmac_test.go
│   │   │   │   ├── linebot_test.go
│   │   │   │   └── linebot.go
│   │   │   ├── linetoken
│   │   │   │   ├── internal
│   │   │   │   │   └── entity.go
│   │   │   │   ├── linetoken_test.go
│   │   │   │   └── linetoken.go
│   │   │   ├── permission
│   │   │   │   ├── internal
│   │   │   │   │   └── entity.go
│   │   │   │   ├── permission_test.go
│   │   │   │   └── permission.go
│   │   │   ├── vc_signal
│   │   │   │   ├── internal
│   │   │   │   │   └── entity.go
│   │   │   │   ├── permission_test.go
│   │   │   │   └── permission.go
│   │   │   └── webhook
│   │   │       ├── internal
│   │   │       │   └── entity.go
│   │   │       ├── webhook_test.go
│   │   │       └── webhook.go
│   │   ├── callback
│   │   │   ├── discord_callback
│   │   │   |   ├── callback_test.go
│   │   │   │   └── callback.go
│   │   │   └── line_callback
│   │   │       ├── callback_test.go
│   │   │       └── callback.go
│   │   ├── login
│   │   │   ├── discord_login
│   │   │   |   ├── discord_login_test.go
│   │   │   │   └── discord_login.go
│   │   │   └── line_login
│   │   │       ├── line_login_test.go
│   │   │       └── line_login.go
│   │   ├── logout
│   │   │   ├── discord_logout
│   │   │   |   ├── discord_logout_test.go
│   │   │   │   └── discord_logout.go
│   │   │   └── line_logout
│   │   │       ├── line_logout_test.go
│   │   │       └── line_logout.go
│   │   └── views
│   │       ├── group
│   │       │   ├── group_test.go
│   │       │   └── group.go
│   │       ├── guildid
│   │       │   ├── line_post_discord_chennel
│   │       │   │   ├── internal
│   │       │   │   │   └── component.go
│   │       │   │   ├── line_post_discord_channel_test.go
│   │       │   │   └── line_post_discord_chennel.go
│   │       │   ├── linetoken
│   │       │   │   ├── internal
│   │       │   │   │   └── component.go
│   │       │   │   ├── linetoken_test.go
│   │       │   │   └── linetoken.go
│   │       │   ├── permission
│   │       │   │   ├── internal
│   │       │   │   │   └── component.go
│   │       │   │   ├── permission_test.go
│   │       │   │   └── permission.go
│   │       │   ├── vc_signal
│   │       │   │   ├── internal
│   │       │   │   │   └── component.go
│   │       │   │   ├── vc_signal_test.go
│   │       │   │   └── vc_signal.go
│   │       │   ├── webhook
│   │       │   │   ├── internal
│   │       │   │   │   └── component.go
│   │       │   │   ├── webhook_test.go
│   │       │   │   └── webhook.go
│   │       │   ├── guildid_test.go
│   │       │   └── guildid.go
│   │       ├── guilds
│   │       │   ├── guilds_test.go
│   │       │   └── guilds.go
│   │       ├── index_test.go
│   │       └── index.go
│   ├── middleware
│   │   ├── discord_oauth_check_test.go
│   │   ├── discord_oauth_check.go
│   │   ├── line_oauth_check_test.go
│   │   ├── line_oauth_check.go
│   │   └── log.go
│   ├── service
│   │   └── index.go
│   ├── shared
│   │   ├── ctxvalue
│   │   │   ├── ctxvalue.go
│   │   │   ├── discordpermissiondata.go
│   │   │   ├── discorduser.go
│   │   │   ├── lineprofile.go
│   │   │   └── lineuser.go
│   │   ├── model
│   │   │   └── entity.go
│   │   └── session
│   │       ├── discord_oauth_token.go
│   │       ├── discord_user.go
│   │       ├── discordstate.go
│   │       ├── guild_id.go
│   │       ├── line_oauth_token.go
│   │       ├── line_user.go
│   │       ├── linenonce.go
│   │       ├── linestate.go
│   │       └── session.go
│   ├── templates
│   │   ├── static
│   │   │   ├── img
│   │   │   │   ├── discord-icon.png
│   │   │   │   ├── line-icon.png
│   │   │   │   ├── ohime.png
│   │   │   │   └── uchuemon.png
│   │   │   └── js
│   │   │       ├── group.js
│   │   │       ├── group.test.js
│   │   │       ├── line_post_discord_chennel.js
│   │   │       ├── line_post_discord_chennel.test.js
│   │   │       ├── linetoken.js
│   │   │       ├── linetoken.test.js
│   │   │       ├── permission.js
│   │   │       ├── permission.test.js
│   │   │       ├── popover.js
│   │   │       ├── vc_signal.js
│   │   │       ├── vc_signal.test.js
│   │   │       ├── webhook.js
│   │   │       └── webhook.test.js
│   │   ├── views
│   │   │   ├── group
│   │   │   │   └── group.html
│   │   │   ├── guildid
│   │   │   │   ├── line_post_discord_chennel.html
│   │   │   │   ├── linetoken.html
│   │   │   │   ├── permission.html
│   │   │   │   ├── vc_signal.html
│   │   │   │   └── webhook.html
│   │   │   ├── guilds
│   │   │   │   └── guilds.html
|   │   │   ├── login
│   │   │   │   └── line_login.html
│   │   │   └── guildid.html
│   │   ├── index.html
│   │   └── layout.html
│   └── main.go
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
|***guild_id***|TEXT|DiscordサーバーのID|
|***type***|TEXT|権限の種類 (line_post_discord_channel, line_bot, vc_signal, webhook)|
|code|BIGINT|Discordの権限コード、詳細は[こちら](https://discord.com/developers/docs/topics/permissions)|

</details>

<details>
    <summary>permissions_user_id</summary>

ユーザーのサーバー権限設定を保存するテーブル
ここに保存されているユーザーは、```permission```と同じ権限を持っているということになる

|カラム名|型|説明|
|---|---|---|
|***guild_id***|TEXT|DiscordサーバーのID|
|***type***|TEXT|権限の種類 (line_post_discord_channel, line_bot, vc_signal, webhook)|
|***user_id***|TEXT|ユーザーID|
|permission|TEXT|権限レベル(read, write, all)|

</details>

<details>
    <summary>permissions_role_id</summary>

ロールのサーバー権限設定を保存するテーブル
ここに保存されているロールは、```permission```と同じ権限を持っているということになる

|カラム名|型|説明|
|---|---|---|
|***guild_id***|TEXT|DiscordサーバーのID|
|***type***|TEXT|権限の種類 (line_post_discord_channel, line_bot, vc_signal, webhook)|
|***role_id***|TEXT|ロールID|
|permission|TEXT|権限レベル(read, write, all)|

</details>

<details>
    <summary>line_post_discord_channel</summary>

DiscordからLINEグループにメッセージを送信するための設定を保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***channel_id***|TEXT|DiscordのチャンネルID|
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
|***channel***|TEXT|DiscordのチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|***type***|INTEGER|メッセージの種類(ピン止め、スレッド、スレッドの返信)|

</details>

<details>
    <summary>line_ng_discord_user_id</summary>

LINEへ送信しないDiscordユーザーを保存するテーブル  
ここに保存されているユーザーを持つユーザーはLINEにメッセージが送信されない  

|カラム名|型|説明|
|---|---|---|
|***channel***|TEXT|DiscordのチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|***user_id***|TEXT|ユーザーID|

</details>

<details>
    <summary>line_ng_discord_role_id</summary>

LINEへ送信しないDiscordロールを保存するテーブル  
ここに保存されているロールを持つユーザーはLINEにメッセージが送信されない  

|カラム名|型|説明|
|---|---|---|
|***channel***|TEXT|DiscordのチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|***role_id***|TEXT|ロールID|

</details>

<details>
    <summary>vc_signal_channel</summary>

ボイスチャンネル入退出の通知設定を保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***vc_channel_id***|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|send_signal|BOOLEAN|L通知を送信するかどうか|
|send_channel_id|TEXT|通知を送信するチャンネルID|
|join_bot|BOOLEAN|ボイスチャンネルにBotが入室したときの通知を送信するかどうか|
|everyone_mention|BOOLEAN|通知を送信するときに@everyoneを使用するかどうか|

</details>

<details>
    <summary>vc_signal_ng_user_id</summary>

指定されたユーザーがボイスチャンネルに参加した場合通知しない

|カラム名|型|説明|
|---|---|---|
|***vc_channel_id***|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|***user_id***|TEXT|ユーザーID|

</details>


<details>
    <summary>vc_signal_ng_role_id</summary>

指定されたロールがボイスチャンネルに参加した場合通知しない

|カラム名|型|説明|
|---|---|---|
|***vc_channel_id***|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|***role_id***|TEXT|ロールID|

</details>

<details>
    <summary>vc_signal_mention_user_id</summary>

ボイスチャンネルの通知の際にメンションするユーザーを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***vc_channel_id***|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|***user_id***|TEXT|ユーザーID|

</details>

<details>
    <summary>vc_signal_mention_role_id</summary>

ボイスチャンネルの通知の際にメンションするロールを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***vc_channel_id***|TEXT|DiscordのボイスチャンネルID|
|guild_id|TEXT|DiscordのサーバーID|
|***role_id***|TEXT|ロールID|

</details>

<details>
    <summary>webhook</summary>

DiscordのWebhookの送信設定を保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***webhook_serial_id***|SERIAL|Webhookの投稿内容の識別ID|
|guild_id|TEXT|DiscordのサーバーID|
|webhook_id|TEXT|WebhookのID|
|subscription_type|TEXT|読み取るもの(YouTube,NicoNico)|
|subscription_id|TEXT|上記のサービスで投稿者を識別できるもの|
|last_posted_at|TIMESTAMP|最後に投稿した日時|

</details>

<details>
    <summary>webhook_user_mention</summary>

Webhookの送信時にメンションするユーザーを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***webhook_serial_id***|SERIAL|Webhookの投稿内容の識別ID|
|***user_id***|TEXT|ユーザーID|

</details>

<details>
    <summary>webhook_role_mention</summary>

Webhookの送信時にメンションするロールを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***webhook_serial_id***|SERIAL|Webhookの投稿内容の識別ID|
|***role_id***|TEXT|ロールID|

</details>

<details>
    <summary>webhook_word</summary>

Webhookの送信時に特定の単語が含まれていた場合にメンションするユーザー、ロールを保存するテーブル  
Twitter運用時使用していたが現在死に要素  
conditionsは投稿時の条件を示す(NgOrはいずれかの単語が含まれていれば投稿しない。SearchAndは全ての単語が含まれていれば投稿する。MentionOrはいずれかの単語が含まれていればメンションする。)

|カラム名|型|説明|
|---|---|---|
|***webhook_serial_id***|SERIAL|Webhookの投稿内容の識別ID|
|***word***|TEXT|メンションする単語|
|conditions|TEXT|投稿時の条件(NgOr NgAnd SearchOr SearchAnd MentionOr MentionAnd)|

</details>

<details>
    <summary>line_bot</summary>

LINEBotの設定を保存するテーブル  
LINEBotのアクセストークン、チャンネルシークレットなどをAES暗号化して保存する

|カラム名|型|説明|
|---|---|---|
|***guild_id***|TEXT|DiscordのサーバーID|
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
    <summary>line_bot_iv</summary>

LINEBotの復号化に使用するIVを保存するテーブル

|カラム名|型|説明|
|---|---|---|
|***guild_id***|TEXT|DiscordのサーバーID|
|line_notify_token_iv|BYTEA|LINE NotifyトークンのIV|
|line_bot_token_iv|BYTEA|LINEBotのアクセストークンのIV|
|line_bot_secret_iv|BYTEA|LINEBotのチャンネルシークレットのIV|
|line_group_id_iv|BYTEA|LINEのグループIDのIV|
|line_client_id_iv|BYTEA|LINEのクライアントIDのIV|
|line_client_secret_iv|BYTEA|LINEのクライアントシークレットのIV|

</details>

# jest
formからjsonに変換する際にjsを使用しているため、jestを使用してテストを行っています。

```bash
npm install
npm run test
```
