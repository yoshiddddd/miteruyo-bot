-- Insert sample data based on ER diagram structure

-- Insert categories
INSERT INTO categories (category_id, category_name) VALUES
    (1, 'General'),
    (2, 'Development'),
    (3, 'Random')
ON CONFLICT (category_id) DO NOTHING;

-- Insert channels
INSERT INTO channels (channel_id, channel_name, category_id) VALUES
    (101, 'general', 1),
    (102, 'announcements', 1),
    (201, 'backend', 2),
    (202, 'frontend', 2),
    (301, 'off-topic', 3)
ON CONFLICT (channel_id) DO NOTHING;

-- Insert users
INSERT INTO users (user_id, user_name) VALUES
    (1001, '田中太郎'),
    (1002, '佐藤花子'),
    (1003, '鈴木次郎'),
    (1004, '高橋美咲'),
    (1005, '山田健太'),
    (1006, '渡辺あゆみ'),
    (1007, '伊藤雄一'),
    (1008, '中村真理'),
    (1009, '小林大輔'),
    (1010, '加藤さくら')
ON CONFLICT (user_id) DO NOTHING;

-- Insert messages
INSERT INTO messages (message_id, user_id, content, channel_id, created_at) VALUES
    (2001, 1001, 'こんにちは！', 101, '2024-01-15 10:30:00'),
    (2002, 1002, 'よろしくお願いします', 101, '2024-01-16 14:22:00'),
    (2003, 1003, 'プロジェクトの進捗はどうですか？', 201, '2024-01-17 09:15:00'),
    (2004, 1004, 'フロントエンドの修正完了しました', 202, '2024-01-18 16:45:00'),
    (2005, 1005, '今日は天気がいいですね', 301, '2024-01-19 11:30:00')
ON CONFLICT (message_id) DO NOTHING;

-- Insert reactions
INSERT INTO reactions (reaction_id, user_id, message_id, emoji, created_at) VALUES
    (3001, 1002, 2001, '👍', '2024-01-15 10:35:00'),
    (3002, 1003, 2001, '😊', '2024-01-15 10:36:00'),
    (3003, 1001, 2002, '🎉', '2024-01-16 14:25:00'),
    (3004, 1005, 2004, '✅', '2024-01-18 16:50:00'),
    (3005, 1006, 2005, '☀️', '2024-01-19 11:35:00')
ON CONFLICT (reaction_id) DO NOTHING;