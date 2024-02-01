package core

import (
	"github.com/bwmarrin/discordgo"
)

func main() {
	Token := "Bot " //"Bot"という接頭辞がないと401 unauthorizedエラーが起きます
	discord, err := discordgo.New(Token)

	// 権限追加
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	discord.Token = Token
	if err != nil {
		panic(err)
	}
	// websocketを開いてlistening開始
	if err = discord.Open(); err != nil {
		panic("Error while opening session")
	}

}