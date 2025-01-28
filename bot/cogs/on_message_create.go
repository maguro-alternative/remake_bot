package cogs

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/pkg/lineworks"

	"github.com/maguro-alternative/remake_bot/bot/cogs/internal"
	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/bot/ffmpeg"

	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
)

func (h *cogHandler) onMessageCreate(s *discordgo.Session, vs *discordgo.MessageCreate) {
	ctx := context.Background()
	repo := repository.NewRepository(h.db)
	ff := ffmpeg.NewFfmpeg(ctx)
	// 暗号化キーのバイトへの変換
	aesCrypto, err := crypto.NewAESCrypto(config.PrivateKey())
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイト変換に失敗しました", "エラー:", err.Error())
	}
	err = onMessageCreateFunc(ctx, h.client, repo, ff, aesCrypto, s, vs)
	if err != nil {
		slog.ErrorContext(ctx, "OnMessageCreate Error", "Error:", err.Error())
	}
	/*err = onMessageCreateFunc2(ctx, h.client, repo, ff, aesCrypto, s, vs)
	if err != nil {
		slog.ErrorContext(ctx, "OnMessageCreate2 Error", "Error:", err.Error())
	}*/
}

func onMessageCreateFunc(
	ctx context.Context,
	client *http.Client,
	repo repository.RepositoryFunc,
	ff ffmpeg.FfmpegInterface,
	aesCrypto crypto.AESInterface,
	s mock.Session,
	vs *discordgo.MessageCreate,
) error {
	var channel repository.LinePostDiscordChannel
	var lineMessageTypes []*line.LineMessageType
	var imageUrls []string
	var videoCount, voiceCount int

	sendTextBuilder := strings.Builder{}

	channel, err := repo.GetLinePostDiscordChannelByChannelID(ctx, vs.ChannelID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.ErrorContext(ctx, "line_post_discord_channelの取得に失敗しました", "エラー:", err.Error())
		return err
	} else if err != nil {
		err = repo.InsertLinePostDiscordChannelByChannelIDAndGuildID(ctx, vs.ChannelID, vs.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_post_discord_channelの登録に失敗しました", "エラー:", err.Error())
			return err
		}
		channel = repository.LinePostDiscordChannel{
			Ng:         false,
			BotMessage: false,
		}
	}
	ngTypes, err := repo.GetLineNgDiscordMessageTypeByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	ngDiscordUserIDs, err := repo.GetLineNgDiscordUserIDByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	ngDiscordRoleIDs, err := repo.GetLineNgDiscordRoleIDByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	// メッセージの種類がNGの場合は処理を終了
	for _, ngType := range ngTypes {
		if vs.Message.Type == discordgo.MessageType(ngType) {
			slog.InfoContext(ctx, "NG Type")
			return err
		}
	}
	// メッセージの送信者がNGの場合は処理を終了
	for _, ngDiscordID := range ngDiscordUserIDs {
		if vs.Author.ID == ngDiscordID {
			slog.InfoContext(ctx, "NG User")
			return err
		}
	}
	for _, ngDiscordRoleID := range ngDiscordRoleIDs {
		for _, role := range vs.Member.Roles {
			if role == ngDiscordRoleID {
				slog.InfoContext(ctx, "NG Role")
				return err
			}
		}
	}
	// チャンネルがNGの場合、またはBotメッセージでない場合は処理を終了
	if channel.Ng || (!channel.BotMessage && vs.Author.Bot) {
		slog.InfoContext(ctx, "NG Channel or Bot Message")
		return err
	}
	lineBotApi, err := repo.GetLineBotNotClientByGuildID(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	lineBotIv, err := repo.GetLineBotIvNotClientByGuildID(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	var lineBotDecrypt internal.LineBotDecrypt

	lineNotifyTokenByte, err := aesCrypto.Decrypt(lineBotApi.LineNotifyToken[0], lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_notify_tokenの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineBotTokenByte, err := aesCrypto.Decrypt(lineBotApi.LineBotToken[0], lineBotIv.LineBotTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_tokenの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineGroupIDByte, err := aesCrypto.Decrypt(lineBotApi.LineGroupID[0], lineBotIv.LineGroupIDIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_group_idの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
	lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
	lineBotDecrypt.LineGroupID = string(lineGroupIDByte)

	lineRequ := line.NewLineRequest(
		*client,
		lineBotDecrypt.LineNotifyToken,
		lineBotDecrypt.LineBotToken,
		lineBotDecrypt.LineGroupID,
	)

	// メッセージの種類によって処理を分岐
	switch vs.Message.Type {
	case discordgo.MessageTypeUserPremiumGuildSubscription:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストしました。")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル1になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル2になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル3になりました！！！！！！！！")
	case discordgo.MessageTypeGuildMemberJoin:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "が参加しました。")
	default:
		st, err := s.Channel(vs.ChannelID)
		if err != nil {
			slog.ErrorContext(ctx, "channel取得に失敗しました", "エラー:", err.Error())
			return err
		}
		sendTextBuilder.WriteString(st.Name + "にて、" + vs.Message.Author.Username)
	}

	// スタンプが送信されていた場合、画像URLを取得
	if vs.StickerItems != nil {
		for _, sticker := range vs.StickerItems {
			switch sticker.FormatType {
			case 1:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".png")
			case 2:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".apng")
			case 4:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".gif")
			}
		}
	}

	videoCount = 0
	voiceCount = 0
	// 添付ファイルが送信されていた場合、LINE用に変換
	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		fileNameNoExt := filepath.Base(attachment.Filename[:len(attachment.Filename)-len(extension)])
		switch extension {
		case ".png", ".jpg", ".jpeg", ".gif":
			imageUrls = append(imageUrls, attachment.URL)
		case ".mp4", ".mov", ".avi", ".wmv", ".flv", ".webm":
			st, err := s.Guild(vs.GuildID)
			if err != nil {
				slog.ErrorContext(ctx, "サーバーの取得に失敗しました", "エラー:", err.Error())
				return err
			}
			lineMessageType := lineRequ.NewLineVideoMessage(attachment.URL, st.IconURL("512"))
			lineMessageTypes = append(lineMessageTypes, lineMessageType)
			videoCount++
		case ".mp3", ".wav", ".ogg", ".m4a":
			tmpFile := os.TempDir() + "/" + attachment.Filename
			tmpFileNotExt := os.TempDir() + "/" + fileNameNoExt
			slog.InfoContext(ctx, "download:"+attachment.URL)
			err = downloadFile(client, tmpFile, attachment.URL)
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのダウンロードに失敗しました", "エラー:", err.Error())
				return err
			}
			if extension != ".m4a" {
				slog.InfoContext(ctx, "m4a変換:"+tmpFile)
				slog.InfoContext(ctx, "ffmpeg:"+tmpFile)
				err = ff.ConversionAudioFile(tmpFile, tmpFileNotExt)
				if err != nil {
					slog.ErrorContext(ctx, "ffmpegの秒数カウントに失敗しました", "エラー:", err.Error())
					return err
				}
			}
			f, err := os.Open(tmpFileNotExt + ".m4a")
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのオープンに失敗しました", "エラー:", err.Error())
				return err
			}
			defer f.Close()
			fileSendMesssage, err := s.ChannelFileSendWithMessage(
				vs.ChannelID,
				"m4aに変換します。",
				tmpFileNotExt+".m4a",
				f,
			)
			if err != nil {
				slog.ErrorContext(ctx, "Discordへのメッセージ送信に失敗しました", "エラー:", err.Error())
				return err
			}
			// 音声ファイルの秒数を取得
			slog.InfoContext(ctx, "秒数取得:"+tmpFileNotExt+".m4a")
			audioLen, err := ff.GetAudioFileSecond(tmpFile, tmpFileNotExt)
			if err != nil {
				slog.ErrorContext(ctx, "音声ファイルの秒数の抽出に失敗しました", "エラー:", err.Error())
				return err
			}
			slog.InfoContext(ctx, "秒数:"+strconv.FormatFloat(audioLen, 'f', -1, 64))
			audio := lineRequ.NewLineAudioMessage(
				fileSendMesssage.Attachments[0].URL,
				audioLen,
			)
			lineMessageTypes = append(lineMessageTypes, audio)
			voiceCount++
		default:
			slog.InfoContext(ctx, "未対応のファイル形式です。")
			sendTextBuilder.WriteString(attachment.URL + "\n")
		}
	}

	if len(imageUrls) > 0 {
		sendTextBuilder.WriteString(" 画像を" + strconv.Itoa(len(imageUrls)) + "枚、")
	}
	if videoCount > 0 {
		sendTextBuilder.WriteString(" 動画を" + strconv.Itoa(videoCount) + "個、")
	}
	if voiceCount > 0 {
		sendTextBuilder.WriteString(" 音声を" + strconv.Itoa(voiceCount) + "個、")
	}
	if len(imageUrls) > 0 || videoCount > 0 || voiceCount > 0 {
		sendTextBuilder.WriteString(" 送信しました。")
	}

	sendTextBuilder.WriteString("「 " + vs.Message.Content + " 」")

	// LINEに送信
	for _, url := range imageUrls {
		err = lineRequ.PushImageNotify(ctx, sendTextBuilder.String(), url)
		if err != nil {
			slog.ErrorContext(ctx, "LINE Notifyの画像送信に失敗しました", "エラー:", err.Error())
			return err
		}
	}
	// 動画、音声を送信
	if len(lineMessageTypes) > 0 {
		err = lineRequ.PushMessageBotInGroup(ctx, lineMessageTypes)
		if err != nil {
			slog.ErrorContext(ctx, "音声、動画の送信に失敗しました", "エラー:", err.Error())
			return err
		}
	}
	// 画像、動画、音声がない場合はテキストのみ送信
	if len(imageUrls) > 0 || videoCount > 0 || voiceCount > 0 {
		return err
	}
	err = lineRequ.PushMessageNotify(ctx, sendTextBuilder.String())
	if err != nil {
		slog.ErrorContext(ctx, "LINE Notifyの送信に失敗しました", "エラー:", err.Error())
		return err
	}
	return err
}

func onMessageCreateFunc2(
	ctx context.Context,
	client *http.Client,
	repo repository.RepositoryFunc,
	ff ffmpeg.FfmpegInterface,
	aesCrypto crypto.AESInterface,
	s mock.Session,
	vs *discordgo.MessageCreate,
) error {
	var channel repository.LinePostDiscordChannel
	var imageUrls []string
	sendTextBuilder := strings.Builder{}

	channel, err := repo.GetLinePostDiscordChannelByChannelID(ctx, vs.ChannelID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.ErrorContext(ctx, "line_post_discord_channelの取得に失敗しました", "エラー:", err.Error())
		return err
	} else if err != nil {
		err = repo.InsertLinePostDiscordChannelByChannelIDAndGuildID(ctx, vs.ChannelID, vs.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_post_discord_channelの登録に失敗しました", "エラー:", err.Error())
			return err
		}
		channel = repository.LinePostDiscordChannel{
			Ng:         false,
			BotMessage: false,
		}
	}
	ngTypes, err := repo.GetLineNgDiscordMessageTypeByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	ngDiscordUserIDs, err := repo.GetLineNgDiscordUserIDByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	ngDiscordRoleIDs, err := repo.GetLineNgDiscordRoleIDByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	// メッセージの種類がNGの場合は処理を終了
	for _, ngType := range ngTypes {
		if vs.Message.Type == discordgo.MessageType(ngType) {
			slog.InfoContext(ctx, "NG Type")
			return err
		}
	}
	// メッセージの送信者がNGの場合は処理を終了
	for _, ngDiscordID := range ngDiscordUserIDs {
		if vs.Author.ID == ngDiscordID {
			slog.InfoContext(ctx, "NG User")
			return err
		}
	}
	for _, ngDiscordRoleID := range ngDiscordRoleIDs {
		for _, role := range vs.Member.Roles {
			if role == ngDiscordRoleID {
				slog.InfoContext(ctx, "NG Role")
				return err
			}
		}
	}
	// チャンネルがNGの場合、またはBotメッセージでない場合は処理を終了
	if channel.Ng || (!channel.BotMessage && vs.Author.Bot) {
		slog.InfoContext(ctx, "NG Channel or Bot Message")
		return err
	}

	lineWorksBotApi, err := repo.GetLineWorksBotByGuildID(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_botの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksBotIv, err := repo.GetLineWorksBotIVByGuildID(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_ivの取得に失敗しました", "エラー:", err.Error())
		return err
	}

	lineWorksBotToken, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksBotToken[0], lineWorksBotIv.LineWorksBotTokenIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_tokenの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksRefreshToken, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksRefreshToken[0], lineWorksBotIv.LineWorksRefreshTokenIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_refresh_tokenの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksGroupID, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksGroupID[0], lineWorksBotIv.LineWorksGroupIDIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_group_idの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksBotID, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksBotID[0], lineWorksBotIv.LineWorksBotIDIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_idの復号化に失敗しました", "エラー:", err.Error())
		return err
	}

	lineworksRequ := lineworks.NewLineWorks(
		*client,
		string(lineWorksBotToken),
		string(lineWorksRefreshToken),
		string(lineWorksBotID),
		string(lineWorksGroupID),
	)

	// メッセージの種類によって処理を分岐
	switch vs.Message.Type {
	case discordgo.MessageTypeUserPremiumGuildSubscription:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストしました。")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル1になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル2になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル3になりました！！！！！！！！")
	case discordgo.MessageTypeGuildMemberJoin:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "が参加しました。")
	default:
		st, err := s.Channel(vs.ChannelID)
		if err != nil {
			slog.ErrorContext(ctx, "channel取得に失敗しました", "エラー:", err.Error())
			return err
		}
		sendTextBuilder.WriteString(st.Name + "にて、" + vs.Message.Author.Username)
	}

	// スタンプが送信されていた場合、画像URLを取得
	if vs.StickerItems != nil {
		for _, sticker := range vs.StickerItems {
			switch sticker.FormatType {
			case 1:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".png")
			case 2:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".apng")
			case 4:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".gif")
			}
		}
	}

	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		switch extension {
		case ".png", ".jpg", ".jpeg", ".gif":
			imageUrls = append(imageUrls, attachment.URL)
		default:
			slog.InfoContext(ctx, "未対応のファイル形式です。")
			sendTextBuilder.WriteString(attachment.URL + "\n")
		}
	}

	if len(imageUrls) > 0 {
		sendTextBuilder.WriteString(" 画像を" + strconv.Itoa(len(imageUrls)) + "枚、送信しました。")
	}
	sendTextBuilder.WriteString("「 " + vs.Message.Content + " 」")

	for _, url := range imageUrls {
		message := lineworksRequ.NewLineWorksImageMessage(url, url)
		lineworksRequ.PushLineWorksMessage(ctx, message)
	}
	message := lineworksRequ.NewLineWorksTextMessage(sendTextBuilder.String())
	res, err := lineworksRequ.PushLineWorksMessage(ctx, message)
	if err != nil {
		return err
	}
	if res.StatusCode == 201 {
		return nil
	}
	slog.InfoContext(ctx, "LINE WORKSのメッセージ送信に失敗しました。トークンを再発行します。", "ステータスコード:", res.StatusCode)
	lineWorksBotInfoApi, err := repo.GetLineWorksBotInfoByGuildID(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_infoの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksBotInfoIv, err := repo.GetLineWorksBotInfoIVByGuildID(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_info_ivの取得に失敗しました", "エラー:", err.Error())
		return err
	}

	lineWorksClientID, err := aesCrypto.Decrypt(lineWorksBotInfoApi.LineWorksClientID[0], lineWorksBotInfoIv.LineWorksClientIDIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_client_idの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksClientSecret, err := aesCrypto.Decrypt(lineWorksBotInfoApi.LineWorksClientSecret[0], lineWorksBotInfoIv.LineWorksClientSecretIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_client_secretの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksServiceAccount, err := aesCrypto.Decrypt(lineWorksBotInfoApi.LineWorksServiceAccount[0], lineWorksBotInfoIv.LineWorksServiceAccountIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_service_accountの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksPrivateKey, err := aesCrypto.Decrypt(lineWorksBotInfoApi.LineWorksPrivateKey[0], lineWorksBotInfoIv.LineWorksPrivateKeyIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_private_keyの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksDomainID, err := aesCrypto.Decrypt(lineWorksBotInfoApi.LineWorksDomainID[0], lineWorksBotInfoIv.LineWorksDomainIDIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_domain_idの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksAdminID, err := aesCrypto.Decrypt(lineWorksBotInfoApi.LineWorksAdminID[0], lineWorksBotInfoIv.LineWorksAdminIDIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_admin_idの復号化に失敗しました", "エラー:", err.Error())
		return err
	}

	lineworksInfoReru := lineworks.NewLineWorksInfo(
		*client,
		string(lineWorksClientID),
		string(lineWorksClientSecret),
		string(lineWorksServiceAccount),
		string(lineWorksPrivateKey),
		string(lineWorksDomainID),
		string(lineWorksAdminID),
	)

	refreshAccessToken, err := lineworksInfoReru.RefreshAccessToken(
		ctx,
		string(lineWorksRefreshToken),
	)
	if err != nil {
		slog.ErrorContext(ctx, "アクセストークンのリフレッシュに失敗しました", "エラー:", err.Error())
		return err
	}

	expiresIn, err := strconv.Atoi(refreshAccessToken.ExpiresIn)
	lineWorksRefreshedTokenExpiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	lineworksRequ = lineworks.NewLineWorks(
		*client,
		refreshAccessToken.AccessToken,
		string(lineWorksRefreshToken),
		string(lineWorksBotID),
		string(lineWorksGroupID),
	)
	res, err = lineworksRequ.PushLineWorksMessage(ctx, message)
	if res.StatusCode == 201 {
		slog.InfoContext(ctx, "LINE WORKSのメッセージ送信に成功しました。")
		lineWorksAccessTokenIv, lineWorksAccessToken, err := aesCrypto.Encrypt(pq.ByteaArray{[]byte(refreshAccessToken.AccessToken)}[0])
		if err != nil {
			slog.ErrorContext(ctx, "line_works_access_tokenの暗号化に失敗しました", "エラー:", err.Error())
			return err
		}
		err = repo.UpdateLineWorksBot(
			ctx,
			repository.NewLineWorksBot(
				vs.GuildID,
				pq.ByteaArray{lineWorksAccessToken},
				lineWorksBotApi.LineWorksRefreshToken,
				lineWorksBotApi.LineWorksGroupID,
				lineWorksBotApi.LineWorksBotID,
				lineWorksBotApi.LineWorksBotSecret,
				pq.NullTime{Time: lineWorksRefreshedTokenExpiresAt, Valid: true},
				lineWorksBotApi.DefaultChannelID,
				lineWorksBotApi.DebugMode,
			),
		)
		if err != nil {
			slog.ErrorContext(ctx, "line_works_botの更新に失敗しました", "エラー:", err.Error())
			return err
		}
		err = repo.UpdateLineWorksBotIV(
			ctx,
			repository.NewLineWorksBotIV(
				vs.GuildID,
				pq.ByteaArray{lineWorksAccessTokenIv},
				lineWorksBotIv.LineWorksRefreshTokenIV,
				lineWorksBotIv.LineWorksGroupIDIV,
				lineWorksBotIv.LineWorksBotIDIV,
				lineWorksBotIv.LineWorksBotSecretIV,
			),
		)
		if err != nil {
			slog.ErrorContext(ctx, "line_works_bot_ivの更新に失敗しました", "エラー:", err.Error())
			return err
		}
		return nil
	}
	if err != nil {
		slog.ErrorContext(ctx, "LINE WORKSのメッセージ送信に失敗しました", "エラー:", err.Error())
		return err
	}

	lineworksTokenInfo, err := lineworksInfoReru.GetAccessToken(
		ctx,
		"bot.message bot.read user.profile.read",
	)
	if err != nil {
		slog.ErrorContext(ctx, "トークンの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksAccessTokenIv, lineWorksAccessToken, err := aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineworksTokenInfo.AccessToken)}[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_access_tokenの暗号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineWorksRefreshTokenIv, lineWorksRefreshToken, err := aesCrypto.Encrypt(pq.ByteaArray{[]byte(lineworksTokenInfo.RefreshToken)}[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_refresh_tokenの暗号化に失敗しました", "エラー:", err.Error())
		return err
	}
	expiresIn, err = strconv.Atoi(lineworksTokenInfo.ExpiresIn)
	lineWorksTokenExpiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	err = repo.UpdateLineWorksBot(
		ctx,
		repository.NewLineWorksBot(
			vs.GuildID,
			pq.ByteaArray{lineWorksAccessToken},
			pq.ByteaArray{lineWorksRefreshToken},
			lineWorksBotApi.LineWorksGroupID,
			lineWorksBotApi.LineWorksBotID,
			lineWorksBotApi.LineWorksBotSecret,
			pq.NullTime{Time: lineWorksTokenExpiresAt, Valid: true},
			lineWorksBotApi.DefaultChannelID,
			lineWorksBotApi.DebugMode,
		),
	)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_botの更新に失敗しました", "エラー:", err.Error())
		return err
	}
	err = repo.UpdateLineWorksBotIV(
		ctx,
		repository.NewLineWorksBotIV(
			vs.GuildID,
			pq.ByteaArray{lineWorksAccessTokenIv},
			pq.ByteaArray{lineWorksRefreshTokenIv},
			lineWorksBotIv.LineWorksGroupIDIV,
			lineWorksBotIv.LineWorksBotIDIV,
			lineWorksBotIv.LineWorksBotSecretIV,
		),
	)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_ivの更新に失敗しました", "エラー:", err.Error())
		return err
	}
	_, err = lineworksRequ.PushLineWorksMessage(ctx, message)
	return err
}

func downloadFile(client *http.Client, tmpFilePath, url string) error {
	f, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
