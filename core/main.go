package main

import (
	"context"
	_ "embed"
	"net"
	"net/http"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/maguro-alternative/remake_bot/bot"
	"github.com/maguro-alternative/remake_bot/core/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web"
	"github.com/maguro-alternative/remake_bot/tasks"

	"github.com/bwmarrin/discordgo"
)

//go:embed schema.sql
var schema string // schema.sqlの内容をschemaに代入
var permissionTypes = []string{"lineBot", "linePostDiscordChannel", "vcSignal", "webhook"}

func main() {
	ctx := context.Background()
	// データベースの接続を開始
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// データベースの初期化
	if _, err := dbV1.ExecContext(ctx, schema); err != nil {
		panic(err)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}

	// ボットの起動
	discord, cleanupCommandHandlers, err := bot.BotOnReady(dbV1, &client)
	if err != nil {
		panic(err)
	}
	defer cleanupCommandHandlers()
	defer discord.Close()

	// データベースにユーザーを追加
	if err := autoDBInsert(ctx, dbV1, discord); err != nil {
		panic(err)
	}

	// サーバーの待ち受けを開始(ゴルーチンで非同期処理)
	// ここでサーバーを起動すると、Ctrl+Cで終了するまでサーバーが起動し続ける
	go func() {
		web.NewWebRouter(
			dbV1,
			&client,
			discord,
		)
	}()

	// webhook投稿のタスクを開始
	go func() {
		tasks.Run(ctx, dbV1, discord)
	}()

	// Ctrl+Cを受け取るためのチャンネル
	sc := make(chan os.Signal, 1)
	// Ctrl+Cを受け取る
	signal.Notify(sc, os.Interrupt)
	<-sc //プログラムが終了しないようロック
}

func autoDBInsert(ctx context.Context, dbv1 db.Driver, discordSession *discordgo.Session) error {
	// データベースにユーザーを追加
	// ここにユーザーを追加するコードを書く
	// 例: dbV1.ExecContext(ctx, "INSERT INTO users (discord_id) VALUES ($1)", discordSession.State.User.ID)
	guilds, err := discordSession.UserGuilds(100, "", "")
	if err != nil {
		return err
	}
	for _, guild := range guilds {
		guildSt, err := discordSession.Guild(guild.ID)
		if err != nil {
			return err
		}
		channels, err := discordSession.GuildChannels(guild.ID)
		if err != nil {
			return err
		}
		slog.InfoContext(ctx, guildSt.Name, "チャンネル数", len(channels))
		channelsAndThreads := make([]*discordgo.Channel, 0, len(channels)+len(guildSt.Threads))
		channelsAndThreads = append(channelsAndThreads, channels...)
		channelsAndThreads = append(channelsAndThreads, guildSt.Threads...)
		activeThreads, err := discordSession.GuildThreadsActive(guild.ID)
		if err != nil {
			return err
		}
		slog.InfoContext(ctx, guildSt.Name, "のアクティブスレッド数", len(activeThreads.Threads))
		channelsAndThreads = append(channelsAndThreads, activeThreads.Threads...)
		for _, channel := range channelsAndThreads {
			if channel.Type == discordgo.ChannelTypeGuildCategory {
				continue
			}
			query := `
				INSERT INTO line_post_discord_channel (
					channel_id,
					guild_id,
					ng,
					bot_message
				) VALUES (
					$1,
					$2,
					$3,
					$4
				) ON CONFLICT (channel_id) DO NOTHING
			`
			if channel.Type == discordgo.ChannelTypeGuildForum {
				var threads *discordgo.ThreadsList
				if len(channel.PermissionOverwrites) == 0 {
					threads, err = discordSession.ThreadsArchived(channel.ID, nil, 100)
					if err != nil {
						return err
					}
				} else {
					threads, err = discordSession.ThreadsPrivateArchived(channel.ID, nil, 100)
					if err != nil {
						return err
					}
				}
				for _, thread := range threads.Threads {
					_, err = dbv1.ExecContext(ctx, query, thread.ID, guild.ID, false, false)
					if err != nil {
						return err
					}
				}
			}
			_, err = dbv1.ExecContext(ctx, query, channel.ID, guild.ID, false, false)
			if err != nil {
				return err
			}
			if channel.Type == discordgo.ChannelTypeGuildVoice {
				query := `
					INSERT INTO vc_signal_channel (
						vc_channel_id,
						guild_id,
						send_signal,
						send_channel_id,
						join_bot,
						everyone_mention
					) VALUES (
						$1,
						$2,
						$3,
						$4,
						$5,
						$6
					) ON CONFLICT (vc_channel_id) DO NOTHING
				`
				_, err := dbv1.ExecContext(ctx, query, channel.ID, guild.ID, true, guildSt.SystemChannelID, false, true)
				if err != nil {
					return err
				}
			}
		}
		for _, permissionType := range permissionTypes {
			query := `
				INSERT INTO permissions_code (
					guild_id,
					type,
					code
				) VALUES (
					$1,
					$2,
					$3
				) ON CONFLICT (guild_id, type) DO NOTHING
			`
			_, err = dbv1.ExecContext(ctx, query, guild.ID, permissionType, 8)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
