package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("DISCORD_TOKEN が設定されていません")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Botの作成に失敗しました:", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Discordへの接続に失敗しました:", err)
		return
	}
	defer dg.Close()

	fmt.Println("Botが起動しました。Ctrl+Cで終了します。")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("Botを終了します。")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	allowedIDs := map[string]bool{
		"1081540011222126705": true,//よっしーさん
		"1132297912207028245": false,//てんとうむし
		"1224248899297087505": false,//たろいもさん
	}

	if allowedIDs[m.Author.ID] {
	s.ChannelMessageSend(m.ChannelID, m.Content)
	}
}