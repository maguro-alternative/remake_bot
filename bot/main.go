package bot

import (
	"context"
	"net/http"
	"strconv"

	"github.com/maguro-alternative/remake_bot/bot/cogs"
	"github.com/maguro-alternative/remake_bot/bot/commands"
	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/pkg/lineworks_service"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
)

func BotOnReady(indexDB db.Driver, client *http.Client) (*discordgo.Session, func(), error) {
	/*
		ボットの起動

		args:
		indexDB: db.Driver
		データベースのドライバー

		return:
		*discordgo.Session
		エラーがなければ、セッションを返します。
		エラーがあれば、エラーを返します。
	*/
	// セッションを作成
	discordToken := "Bot " + config.DiscordBotToken()
	discordSession, err := discordgo.New(discordToken)
	if err != nil {
		return nil, func() {}, errors.WithStack(err)
	}

	// LINE Works 設定を環境変数から取得
	worksID := config.LineWorksID()
	password := config.LineWorksPassword()
	channelNoStr := config.ChannelNo()
	
	var channelNo int
	if channelNoStr != "" {
		var err error
		channelNo, err = strconv.Atoi(channelNoStr)
		if err != nil {
			channelNo = 0 // Invalid channel number
		}
	}

	// LINE Works サービス作成（統一設定）
	lineWorksService, err := lineworks_service.NewService(worksID, password, channelNo)
	if err != nil {
		return nil, func() {}, errors.WithStack(err)
	}

	// バックグラウンドサービス開始
	ctx := context.Background()
	if err := lineWorksService.Start(ctx); err != nil {
		return nil, func() {}, errors.WithStack(err)
	}

	discordSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	discordSession.Token = discordToken
	err = discordSession.Open()
	if err != nil {
		lineWorksService.Stop()
		return nil, func() {}, errors.WithStack(err)
	}

	cogs.RegisterHandlers(discordSession, indexDB, client, lineWorksService)
	cleanupCommandHandlers, err := commands.RegisterCommands(discordSession, indexDB, client)

	// クリーンアップ関数に LINE Works サービス停止を追加
	cleanup := func() {
		lineWorksService.Stop()
		cleanupCommandHandlers()
	}

	return discordSession, cleanup, err
}
