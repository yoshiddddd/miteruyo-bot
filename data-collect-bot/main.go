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

	// DB 接続確認
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
	// Botの発言は無視
	if m.Author.Bot {
		return
	}

	userID := m.Author.ID
	userName := m.Author.Username

	// 重複防止：すでに存在する場合は INSERT しない
	_, err := db.Exec(`
		INSERT INTO users (user_id, user_name)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING
	`, userID, userName)
	if err != nil {
		log.Println("DB insert error:", err)
	} else {
		log.Printf("Inserted user %s (%s)\n", userName, userID)
	}
}
