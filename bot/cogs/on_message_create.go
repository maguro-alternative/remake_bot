package cogs

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
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
	/*err = onMessageCreateFunc(ctx, h.client, repo, ff, aesCrypto, s, vs)
	if err != nil {
		slog.ErrorContext(ctx, "OnMessageCreate Error", "Error:", err.Error())
	}*/
	err = onMessageCreateFunc2(ctx, h.client, repo, ff, aesCrypto, s, vs)
	if err != nil {
		slog.ErrorContext(ctx, "OnMessageCreate2 Error", "Error:", err.Error())
	}
	err = onMessageCreateFunc3(ctx, h.client, repo, ff, aesCrypto, s, vs)
	if err != nil {
		slog.ErrorContext(ctx, "OnMessageCreate Error", "Error:", err.Error())
	}
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
	// 1. チャンネル情報取得
	channel, err := getOrCreateChannel(ctx, repo, vs)
	if err != nil {
		return err
	}

	// 2. 権限チェック
	err = validateMessagePermissions(ctx, repo, vs, channel)
	if err != nil {
		return err
	}

	// 3. LINE Bot認証情報取得
	lineBotDecrypt, err := getDecryptedLineBotCredentials(ctx, repo, aesCrypto, vs.GuildID)
	if err != nil {
		return err
	}

	// 4. LINE APIクライアント作成
	lineRequ := line.NewLineRequest(
		*client,
		lineBotDecrypt.LineNotifyToken,
		lineBotDecrypt.LineBotToken,
		lineBotDecrypt.LineGroupID,
	)

	// 5. メッセージテキスト生成
	messageText, err := buildMessageText(ctx, s, vs)
	if err != nil {
		return err
	}
	var sendTextBuilder strings.Builder
	sendTextBuilder.WriteString(messageText)

	// 6. スタンプ処理
	stickerImageUrls := processStickerItems(vs)

	// 7. 添付ファイル処理
	attachmentResult, err := processAttachmentsForLine(ctx, client, s, ff, lineRequ, vs)
	if err != nil {
		return err
	}

	// 8. 全画像URLをマージ
	allImageUrls := append(stickerImageUrls, attachmentResult.ImageUrls...)

	// 9. 未対応ファイルURLを追加
	for _, url := range attachmentResult.UnsupportedUrls {
		sendTextBuilder.WriteString(url + "\n")
	}

	// 10. 統計情報追加
	videoCount := len(attachmentResult.LineMessageTypes)
	voiceCount := 0
	for _, msg := range attachmentResult.LineMessageTypes {
		if msg.Type == "audio" {
			voiceCount++
		}
	}
	if len(allImageUrls) > 0 {
		sendTextBuilder.WriteString(" 画像を" + strconv.Itoa(len(allImageUrls)) + "枚、")
	}
	if videoCount > 0 {
		sendTextBuilder.WriteString(" 動画を" + strconv.Itoa(videoCount-voiceCount) + "個、")
	}
	if voiceCount > 0 {
		sendTextBuilder.WriteString(" 音声を" + strconv.Itoa(voiceCount) + "個、")
	}
	if len(allImageUrls) > 0 || videoCount > 0 {
		sendTextBuilder.WriteString(" 送信しました。")
	}
	sendTextBuilder.WriteString("「 " + vs.Message.Content + " 」")

	finalText := sendTextBuilder.String()

	// 11. LINEに送信
	for _, url := range allImageUrls {
		err = lineRequ.PushImageNotify(ctx, finalText, url)
		if err != nil {
			slog.ErrorContext(ctx, "LINE Notifyの画像送信に失敗しました", "エラー:", err.Error())
			return err
		}
	}

	if len(attachmentResult.LineMessageTypes) > 0 {
		err = lineRequ.PushMessageBotInGroup(ctx, attachmentResult.LineMessageTypes)
		if err != nil {
			slog.ErrorContext(ctx, "音声、動画の送信に失敗しました", "エラー:", err.Error())
			return err
		}
	}

	if len(allImageUrls) == 0 && len(attachmentResult.LineMessageTypes) == 0 {
		err = lineRequ.PushMessageNotify(ctx, finalText)
		if err != nil {
			slog.ErrorContext(ctx, "LINE Notifyの送信に失敗しました", "エラー:", err.Error())
			return err
		}
	}

	return nil
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
	// 1. チャンネル情報取得
	channel, err := getOrCreateChannel(ctx, repo, vs)
	if err != nil {
		return err
	}

	// 2. 権限チェック
	err = validateMessagePermissions(ctx, repo, vs, channel)
	if err != nil {
		return err
	}

	// 3. LINE WORKS認証情報取得
	lineWorksCredentials, err := getDecryptedLineWorksCredentials(ctx, repo, aesCrypto, vs.GuildID)
	if err != nil {
		return err
	}

	// 取得したデータを後続の処理で使うため、変数に保持
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

	// 4. LINE WORKSクライアント作成
	lineworksRequ := lineworks.NewLineWorks(
		*client,
		lineWorksCredentials.BotToken,
		lineWorksCredentials.RefreshToken,
		lineWorksCredentials.BotID,
		lineWorksCredentials.GroupID,
	)

	// 5. メッセージテキスト生成
	messageText, err := buildMessageText(ctx, s, vs)
	if err != nil {
		return err
	}
	var sendTextBuilder strings.Builder
	sendTextBuilder.WriteString(messageText)

	// 6. スタンプ処理
	stickerImageUrls := processStickerItems(vs)

	// 7. 画像添付ファイル処理
	attachmentImageUrls := processImageAttachments(vs)

	// 8. 全画像URLをマージ
	allImageUrls := append(stickerImageUrls, attachmentImageUrls...)

	// 9. 統計情報追加
	if len(allImageUrls) > 0 {
		sendTextBuilder.WriteString(" 画像を" + strconv.Itoa(len(allImageUrls)) + "枚、送信しました。")
	}
	sendTextBuilder.WriteString("「 " + vs.Message.Content + " 」")

	// 10. LINE WORKSに画像送信
	for _, url := range allImageUrls {
		message := lineworksRequ.NewLineWorksImageMessage(url, url)
		lineworksRequ.PushLineWorksMessage(ctx, message)
	}

	// 11. テキストメッセージ送信
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
		lineWorksCredentials.RefreshToken,
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
		lineWorksCredentials.RefreshToken,
		lineWorksCredentials.BotID,
		lineWorksCredentials.GroupID,
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

func onMessageCreateFunc3(
	ctx context.Context,
	client *http.Client,
	repo repository.RepositoryFunc,
	ff ffmpeg.FfmpegInterface,
	aesCrypto crypto.AESInterface,
	s mock.Session,
	vs *discordgo.MessageCreate,
) error {
	// 1. チャンネル情報取得
	channel, err := getOrCreateChannel(ctx, repo, vs)
	if err != nil {
		return err
	}

	// 2. 権限チェック
	err = validateMessagePermissions(ctx, repo, vs, channel)
	if err != nil {
		return err
	}

	// 3. メッセージテキスト生成
	messageText, err := buildMessageText(ctx, s, vs)
	if err != nil {
		return err
	}
	var sendTextBuilder strings.Builder
	sendTextBuilder.WriteString(messageText)

	// 4. スタンプ処理
	stickerImageUrls := processStickerItems(vs)

	// 5. 添付ファイル処理（URLのみ取得）
	imageUrls, videoUrls, voiceUrls, unsupportedUrls := processAttachmentsSimple(vs)

	// 6. 全画像URLをマージ
	allImageUrls := append(stickerImageUrls, imageUrls...)

	// 7. 未対応ファイルURLを追加
	for _, url := range unsupportedUrls {
		sendTextBuilder.WriteString(url + "\n")
	}

	// 8. 統計情報追加
	if len(allImageUrls) > 0 {
		sendTextBuilder.WriteString(" 画像を" + strconv.Itoa(len(allImageUrls)) + "枚、")
	}
	if len(videoUrls) > 0 {
		sendTextBuilder.WriteString(" 動画を" + strconv.Itoa(len(videoUrls)) + "個、")
	}
	if len(voiceUrls) > 0 {
		sendTextBuilder.WriteString(" 音声を" + strconv.Itoa(len(voiceUrls)) + "個、")
	}
	if len(allImageUrls) > 0 || len(videoUrls) > 0 || len(voiceUrls) > 0 {
		sendTextBuilder.WriteString(" 送信しました。")
	}
	sendTextBuilder.WriteString("「 " + vs.Message.Content + " 」")

	// 9. 全メディアURLを追加
	for _, url := range allImageUrls {
		sendTextBuilder.WriteString(url + "\n")
	}
	for _, url := range videoUrls {
		sendTextBuilder.WriteString(url + "\n")
	}
	for _, url := range voiceUrls {
		sendTextBuilder.WriteString(url + "\n")
	}

	// 10. 内部APIに送信
	form := url.Values{}
	form.Add("message", sendTextBuilder.String())
	req, err := http.NewRequest(http.MethodPost, config.InternalURL(), strings.NewReader(form.Encode()))
	if err != nil {
		slog.ErrorContext(ctx, "リクエストの作成に失敗しました", "エラー:", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+config.ChannelNo())
	resp, err := client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "リクエストの送信に失敗しました", "エラー:", err.Error())
		return err
	}
	defer resp.Body.Close()

	return nil
}

// getOrCreateChannel チャンネル情報を取得、存在しない場合は作成
func getOrCreateChannel(
	ctx context.Context,
	repo repository.RepositoryFunc,
	vs *discordgo.MessageCreate,
) (repository.LinePostDiscordChannel, error) {
	channel, err := repo.GetLinePostDiscordChannelByChannelID(ctx, vs.ChannelID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.ErrorContext(ctx, "line_post_discord_channelの取得に失敗しました", "エラー:", err.Error())
		return repository.LinePostDiscordChannel{}, err
	} else if err != nil {
		err = repo.InsertLinePostDiscordChannelByChannelIDAndGuildID(ctx, vs.ChannelID, vs.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_post_discord_channelの登録に失敗しました", "エラー:", err.Error())
			return repository.LinePostDiscordChannel{}, err
		}
		return repository.LinePostDiscordChannel{
			Ng:         false,
			BotMessage: false,
		}, nil
	}
	return channel, nil
}

// validateMessagePermissions メッセージの権限をチェック（NG判定）
func validateMessagePermissions(
	ctx context.Context,
	repo repository.RepositoryFunc,
	vs *discordgo.MessageCreate,
	channel repository.LinePostDiscordChannel,
) error {
	// チャンネルがNGの場合、またはBotメッセージでない場合は処理を終了
	if channel.Ng || (!channel.BotMessage && vs.Author.Bot) {
		slog.InfoContext(ctx, "NG Channel or Bot Message")
		return nil
	}

	// NGメッセージタイプのチェック
	ngTypes, err := repo.GetLineNgDiscordMessageTypeByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	for _, ngType := range ngTypes {
		if vs.Message.Type == discordgo.MessageType(ngType) {
			slog.InfoContext(ctx, "NG Type")
			return nil
		}
	}

	// NGユーザーのチェック
	ngDiscordUserIDs, err := repo.GetLineNgDiscordUserIDByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_user_idの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	for _, ngDiscordID := range ngDiscordUserIDs {
		if vs.Author.ID == ngDiscordID {
			slog.InfoContext(ctx, "NG User")
			return nil
		}
	}

	// NGロールのチェック
	ngDiscordRoleIDs, err := repo.GetLineNgDiscordRoleIDByChannelID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_role_idの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	for _, ngDiscordRoleID := range ngDiscordRoleIDs {
		for _, role := range vs.Member.Roles {
			if role == ngDiscordRoleID {
				slog.InfoContext(ctx, "NG Role")
				return nil
			}
		}
	}

	return nil
}

// buildMessageText メッセージの種類によってテキストを生成
func buildMessageText(
	ctx context.Context,
	s mock.Session,
	vs *discordgo.MessageCreate,
) (string, error) {
	var builder strings.Builder

	switch vs.Message.Type {
	case discordgo.MessageTypeUserPremiumGuildSubscription:
		builder.WriteString(vs.Message.Author.Username + "がサーバーブーストしました。")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne:
		builder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル1になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo:
		builder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル2になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree:
		builder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル3になりました！！！！！！！！")
	case discordgo.MessageTypeGuildMemberJoin:
		builder.WriteString(vs.Message.Author.Username + "が参加しました。")
	default:
		st, err := s.Channel(vs.ChannelID)
		if err != nil {
			slog.ErrorContext(ctx, "channel取得に失敗しました", "エラー:", err.Error())
			return "", err
		}
		builder.WriteString(st.Name + "にて、" + vs.Message.Author.Username)
	}

	return builder.String(), nil
}

// processStickerItems スタンプから画像URLを取得
func processStickerItems(vs *discordgo.MessageCreate) []string {
	var imageUrls []string
	if vs.StickerItems == nil {
		return imageUrls
	}

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
	return imageUrls
}

// getDecryptedLineBotCredentials LINE Botの認証情報を取得・復号化
func getDecryptedLineBotCredentials(
	ctx context.Context,
	repo repository.RepositoryFunc,
	aesCrypto crypto.AESInterface,
	guildID string,
) (*internal.LineBotDecrypt, error) {
	lineBotApi, err := repo.GetLineBotNotClientByGuildID(ctx, guildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineBotIv, err := repo.GetLineBotIvNotClientByGuildID(ctx, guildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました", "エラー:", err.Error())
		return nil, err
	}

	var lineBotDecrypt internal.LineBotDecrypt

	lineNotifyTokenByte, err := aesCrypto.Decrypt(lineBotApi.LineNotifyToken[0], lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_notify_tokenの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineBotTokenByte, err := aesCrypto.Decrypt(lineBotApi.LineBotToken[0], lineBotIv.LineBotTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_tokenの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineGroupIDByte, err := aesCrypto.Decrypt(lineBotApi.LineGroupID[0], lineBotIv.LineGroupIDIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_group_idの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}

	lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
	lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
	lineBotDecrypt.LineGroupID = string(lineGroupIDByte)

	return &lineBotDecrypt, nil
}

// LineWorksCredentials LINE WORKS認証情報
type LineWorksCredentials struct {
	BotToken     string
	RefreshToken string
	GroupID      string
	BotID        string
	BotSecret    string
}

// getDecryptedLineWorksCredentials LINE WORKSの認証情報を取得・復号化
func getDecryptedLineWorksCredentials(
	ctx context.Context,
	repo repository.RepositoryFunc,
	aesCrypto crypto.AESInterface,
	guildID string,
) (*LineWorksCredentials, error) {
	lineWorksBotApi, err := repo.GetLineWorksBotByGuildID(ctx, guildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_botの取得に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineWorksBotIv, err := repo.GetLineWorksBotIVByGuildID(ctx, guildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_ivの取得に失敗しました", "エラー:", err.Error())
		return nil, err
	}

	lineWorksBotToken, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksBotToken[0], lineWorksBotIv.LineWorksBotTokenIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_tokenの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineWorksRefreshToken, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksRefreshToken[0], lineWorksBotIv.LineWorksRefreshTokenIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_refresh_tokenの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineWorksGroupID, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksGroupID[0], lineWorksBotIv.LineWorksGroupIDIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_group_idの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineWorksBotID, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksBotID[0], lineWorksBotIv.LineWorksBotIDIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_idの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}
	lineWorksBotSecret, err := aesCrypto.Decrypt(lineWorksBotApi.LineWorksBotSecret[0], lineWorksBotIv.LineWorksBotSecretIV[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_works_bot_secretの復号化に失敗しました", "エラー:", err.Error())
		return nil, err
	}

	return &LineWorksCredentials{
		BotToken:     string(lineWorksBotToken),
		RefreshToken: string(lineWorksRefreshToken),
		GroupID:      string(lineWorksGroupID),
		BotID:        string(lineWorksBotID),
		BotSecret:    string(lineWorksBotSecret),
	}, nil
}

// AttachmentResult 添付ファイル処理の結果
type AttachmentResult struct {
	ImageUrls        []string
	VideoUrls        []string
	VoiceUrls        []string
	LineMessageTypes []*line.LineMessageType
	UnsupportedUrls  []string
}

// processImageAttachments 画像添付ファイルの処理
func processImageAttachments(vs *discordgo.MessageCreate) []string {
	var imageUrls []string
	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		switch extension {
		case ".png", ".jpg", ".jpeg", ".gif":
			imageUrls = append(imageUrls, attachment.URL)
		}
	}
	return imageUrls
}

// processAttachmentsForLine LINE用に添付ファイルを処理（動画・音声含む）
func processAttachmentsForLine(
	ctx context.Context,
	client *http.Client,
	s mock.Session,
	ff ffmpeg.FfmpegInterface,
	lineRequ *line.LineRequest,
	vs *discordgo.MessageCreate,
) (*AttachmentResult, error) {
	result := &AttachmentResult{
		ImageUrls:        []string{},
		LineMessageTypes: []*line.LineMessageType{},
		UnsupportedUrls:  []string{},
	}

	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		fileNameNoExt := filepath.Base(attachment.Filename[:len(attachment.Filename)-len(extension)])

		switch extension {
		case ".png", ".jpg", ".jpeg", ".gif":
			result.ImageUrls = append(result.ImageUrls, attachment.URL)

		case ".mp4", ".mov", ".avi", ".wmv", ".flv", ".webm":
			st, err := s.Guild(vs.GuildID)
			if err != nil {
				slog.ErrorContext(ctx, "サーバーの取得に失敗しました", "エラー:", err.Error())
				return nil, err
			}
			lineMessageType := lineRequ.NewLineVideoMessage(attachment.URL, st.IconURL("512"))
			result.LineMessageTypes = append(result.LineMessageTypes, lineMessageType)

		case ".mp3", ".wav", ".ogg", ".m4a":
			tmpFile := os.TempDir() + "/" + attachment.Filename
			tmpFileNotExt := os.TempDir() + "/" + fileNameNoExt
			slog.InfoContext(ctx, "download:"+attachment.URL)

			err := downloadFile(client, tmpFile, attachment.URL)
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのダウンロードに失敗しました", "エラー:", err.Error())
				return nil, err
			}

			if extension != ".m4a" {
				slog.InfoContext(ctx, "m4a変換:"+tmpFile)
				err = ff.ConversionAudioFile(tmpFile, tmpFileNotExt)
				if err != nil {
					slog.ErrorContext(ctx, "ffmpegの変換に失敗しました", "エラー:", err.Error())
					return nil, err
				}
			}

			f, err := os.Open(tmpFileNotExt + ".m4a")
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのオープンに失敗しました", "エラー:", err.Error())
				return nil, err
			}
			defer f.Close()

			fileSendMessage, err := s.ChannelFileSendWithMessage(
				vs.ChannelID,
				"m4aに変換します。",
				tmpFileNotExt+".m4a",
				f,
			)
			if err != nil {
				slog.ErrorContext(ctx, "Discordへのメッセージ送信に失敗しました", "エラー:", err.Error())
				return nil, err
			}

			slog.InfoContext(ctx, "秒数取得:"+tmpFileNotExt+".m4a")
			audioLen, err := ff.GetAudioFileSecond(tmpFile, tmpFileNotExt)
			if err != nil {
				slog.ErrorContext(ctx, "音声ファイルの秒数の抽出に失敗しました", "エラー:", err.Error())
				return nil, err
			}
			slog.InfoContext(ctx, "秒数:"+strconv.FormatFloat(audioLen, 'f', -1, 64))

			audio := lineRequ.NewLineAudioMessage(fileSendMessage.Attachments[0].URL, audioLen)
			result.LineMessageTypes = append(result.LineMessageTypes, audio)

		default:
			slog.InfoContext(ctx, "未対応のファイル形式です。")
			result.UnsupportedUrls = append(result.UnsupportedUrls, attachment.URL)
		}
	}

	return result, nil
}

// processAttachmentsSimple シンプルな添付ファイル処理（URL取得のみ）
func processAttachmentsSimple(vs *discordgo.MessageCreate) (imageUrls, videoUrls, voiceUrls, unsupportedUrls []string) {
	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		switch extension {
		case ".png", ".jpg", ".jpeg", ".gif":
			imageUrls = append(imageUrls, attachment.URL)
		case ".mp4", ".mov", ".avi", ".wmv", ".flv", ".webm":
			videoUrls = append(videoUrls, attachment.URL)
		case ".mp3", ".wav", ".ogg", ".m4a":
			voiceUrls = append(voiceUrls, attachment.URL)
		default:
			slog.InfoContext(context.Background(), "未対応のファイル形式です。")
			unsupportedUrls = append(unsupportedUrls, attachment.URL)
		}
	}
	return
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
