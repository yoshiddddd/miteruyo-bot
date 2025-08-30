package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "github.com/lib/pq"
	"github.com/bwmarrin/discordgo"
)

var db *sql.DB

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	// DB 接続
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// Discord セッション開始
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	// イベントハンドラを登録
	dg.AddHandler(messageCreate)
	dg.AddHandler(reactionAdd)

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// 終了待ち
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	dg.Close()
	db.Close()
}

// メッセージ受信時の処理
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	userID := m.Author.ID
	userName := m.Author.Username
	channelID := m.ChannelID
	content := m.Content

	// ユーザーを保存
	_, err := db.Exec(`
		INSERT INTO users (user_id, user_name)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING
	`, userID, userName)
	if err != nil {
		log.Println("DB insert error (users):", err)
	}

	// チャンネル情報を取得
	ch, err := s.Channel(channelID)
	if err != nil {
		log.Println("Discord channel fetch error:", err)
		return
	}

	// カテゴリを保存
	categoryID := ch.ParentID
	if categoryID != "" {
		_, err = db.Exec(`
			INSERT INTO categories (category_id, category_name)
			VALUES ($1, $2)
			ON CONFLICT (category_id) DO NOTHING
		`, categoryID, ch.GuildID) // 名前はGuildIDでもよい。必要ならch.ParentIDの名前取得
		if err != nil {
			log.Println("DB insert error (categories):", err)
		}
	}

	// チャンネルを保存
	_, err = db.Exec(`
		INSERT INTO channels (channel_id, channel_name, category_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (channel_id) DO NOTHING
	`, channelID, ch.Name, categoryID)
	if err != nil {
		log.Println("DB insert error (channels):", err)
	}

	// メッセージを保存
	_, err = db.Exec(`
		INSERT INTO messages (message_id, user_id, channel_id, content)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (message_id) DO NOTHING
	`, m.ID, userID, channelID, content)
	if err != nil {
		log.Println("DB insert error (messages):", err)
	} else {
		log.Printf("Inserted message %s from user %s\n", m.ID, userName)
	}
}

// リアクション追加時の処理
func reactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	// リアクションを保存
	_, err := db.Exec(`
		INSERT INTO reactions (reaction_id, user_id, message_id, emoji)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (reaction_id) DO NOTHING
	`, fmt.Sprintf("%s_%s", r.MessageID, r.Emoji.Name), r.UserID, r.MessageID, r.Emoji.Name)
	if err != nil {
		log.Println("DB insert error (reactions):", err)
	} else {
		log.Printf("Inserted reaction %s by user %s\n", r.Emoji.Name, r.UserID)
	}
}