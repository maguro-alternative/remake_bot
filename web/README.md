# web
## ディレクトリ構成
```
├── web
│   ├── components
│   │   ├── channel_select.go
│   │   ├── discord_account_pop.go
│   │   ├── entity.go
│   │   ├── line_account_pop.go
│   │   ├── line_post_discord_chennel.go
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
│   │   │   │   └── linetoken.html
│   │   │   ├── guilds
│   │   │   │   └── guilds.html
|   │   │   ├── login
│   │   │   │   └── line_login.html
│   │   │   └── guildid.html
│   │   ├── index.html
│   │   └── layout.html
│   └── main.go
```

## components
HTMLのコンポーネントを管理するディレクトリです。

## config
環境変数を読み込みます。
```web```内ではここから読み込んだ環境変数を参照します。

## handler
APIのエンドポイントを管理するディレクトリです。
```views```ディレクトリ内ではHTMLファイルを表示するための処理を行います。

## middleware
ミドルウェアを管理するディレクトリです。

## service
プロパティを管理するディレクトリです。

## shared
```web```共通の処理を管理するディレクトリです。

## templates
HTMLテンプレートを管理するディレクトリです。

## main.go
webのメインとなるコードです。
