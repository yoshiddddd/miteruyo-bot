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
	 "github.com/joho/godotenv"
)

var db *sql.DB

// コールドスタート時に全ユーザーを取得する関数
func fetchAllUsers(s *discordgo.Session) error {
	log.Println("Starting to fetch all users from all guilds...")

	// ボットが参加している全ギルドを取得
	guilds := s.State.Guilds

	for _, guild := range guilds {
		log.Printf("Fetching users from guild: %s (%s)", guild.Name, guild.ID)

		// ギルドメンバーを取得（1000件ずつ）
		after := ""
		for {
			members, err := s.GuildMembers(guild.ID, after, 1000)
			if err != nil {
				log.Printf("Error fetching members from guild %s: %v", guild.ID, err)
				break
			}

			if len(members) == 0 {
				break
			}

			// メンバーをDBに保存
			for _, member := range members {
				if member.User.Bot {
					continue // ボットは除外
				}

				// 表示名を決定
				displayName := member.User.Username
				if member.Nick != "" {
					displayName = member.Nick
				}

				// ユーザーをDBに保存
				_, err := db.Exec(`
					INSERT INTO users (user_id, username, display_name)
					VALUES ($1, $2, $3)
					ON CONFLICT (user_id) DO UPDATE SET
						username = EXCLUDED.username,
						display_name = EXCLUDED.display_name,
						updated_at = CURRENT_TIMESTAMP
				`, member.User.ID, member.User.Username, displayName)

				if err != nil {
					log.Printf("DB insert error for user %s: %v", member.User.Username, err)
				}
			}

			log.Printf("Processed %d members from guild %s", len(members), guild.Name)

			// 次のページがある場合
			if len(members) < 1000 {
				break
			}
			after = members[len(members)-1].User.ID
		}
	}

	log.Println("Finished fetching all users")
	return nil
}

func main() {
	errs := godotenv.Load()
	if errs != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set.")
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

	// コールドスタート時に全ユーザーを取得
	err = fetchAllUsers(dg)
	if err != nil {
		log.Println("Warning: Failed to fetch all users on startup:", err)
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

	// 表示名を取得（ニックネームがあれば優先、なければユーザー名）
	displayName := userName
	if m.Member != nil && m.Member.Nick != "" {
		displayName = m.Member.Nick
	}

	// ユーザーを保存
	_, err := db.Exec(`
		INSERT INTO users (user_id, username, display_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			username = EXCLUDED.username,
			display_name = EXCLUDED.display_name,
			updated_at = CURRENT_TIMESTAMP
	`, userID, userName, displayName)
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
		// カテゴリ情報を取得
		category, err := s.Channel(categoryID)
		var categoryName string
		if err != nil {
			log.Println("Failed to get category info, using ID as name:", err)
			categoryName = categoryID // フォールバック：IDを使用
		} else {
			categoryName = category.Name
		}

		_, err = db.Exec(`
			INSERT INTO categories (category_id, category_name)
			VALUES ($1, $2)
			ON CONFLICT (category_id) DO NOTHING
		`, categoryID, categoryName)
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

	// リアクションユーザーの情報を取得（ギルドメンバーから）
	member, err := s.GuildMember(r.GuildID, r.UserID)
	if err != nil {
		log.Println("Failed to get guild member for reaction:", err)
		return
	}

	// ユーザーを保存（存在しない場合）
	var displayName *string
	if member.Nick != "" {
		displayName = &member.Nick
	}
	
	_, err = db.Exec(`
		INSERT INTO users (user_id, username, display_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			username = EXCLUDED.username,
			display_name = EXCLUDED.display_name,
			updated_at = CURRENT_TIMESTAMP
	`, r.UserID, member.User.Username, displayName)
	if err != nil {
		log.Println("DB insert error (users in reaction):", err)
		return
	}

	// メッセージを取得してDBに保存（存在しない場合）
	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Println("Failed to get message for reaction:", err)
		return
	}

	// チャンネル情報を取得
	ch, err := s.Channel(r.ChannelID)
	if err != nil {
		log.Println("Failed to get channel for reaction:", err)
		return
	}

	// カテゴリを保存
	categoryID := ch.ParentID
	if categoryID != "" {
		// カテゴリ情報を取得
		category, err := s.Channel(categoryID)
		var categoryName string
		if err != nil {
			log.Println("Failed to get category info in reaction, using ID as name:", err)
			categoryName = categoryID // フォールバック：IDを使用
		} else {
			categoryName = category.Name
		}

		_, err = db.Exec(`
			INSERT INTO categories (category_id, category_name)
			VALUES ($1, $2)
			ON CONFLICT (category_id) DO NOTHING
		`, categoryID, categoryName)
		if err != nil {
			log.Println("DB insert error (categories in reaction):", err)
		}
	}

	// チャンネルを保存
	_, err = db.Exec(`
		INSERT INTO channels (channel_id, channel_name, category_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (channel_id) DO NOTHING
	`, r.ChannelID, ch.Name, categoryID)
	if err != nil {
		log.Println("DB insert error (channels in reaction):", err)
	}

	// メッセージ作成者のユーザー情報も保存
	_, err = db.Exec(`
		INSERT INTO users (user_id, username, display_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			username = EXCLUDED.username,
			display_name = EXCLUDED.display_name,
			updated_at = CURRENT_TIMESTAMP
	`, msg.Author.ID, msg.Author.Username, nil) // GlobalNameはGuildMemberから取得する必要があります
	if err != nil {
		log.Println("DB insert error (message author in reaction):", err)
	}

	// メッセージを保存
	_, err = db.Exec(`
		INSERT INTO messages (message_id, user_id, channel_id, content)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (message_id) DO NOTHING
	`, r.MessageID, msg.Author.ID, r.ChannelID, msg.Content)
	if err != nil {
		log.Println("DB insert error (messages in reaction):", err)
	}

	// リアクションを保存
	_, err = db.Exec(`
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