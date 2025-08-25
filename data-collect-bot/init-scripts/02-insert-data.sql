-- Insert sample data from CSV
INSERT INTO users (user_name, email, created_at, status, score) VALUES
    ('田中太郎', 'tanaka@example.com', '2024-01-15 10:30:00', 'active', 85),
    ('佐藤花子', 'sato@example.com', '2024-01-16 14:22:00', 'active', 92),
    ('鈴木次郎', 'suzuki@example.com', '2024-01-17 09:15:00', 'inactive', 67),
    ('高橋美咲', 'takahashi@example.com', '2024-01-18 16:45:00', 'active', 78),
    ('山田健太', 'yamada@example.com', '2024-01-19 11:30:00', 'pending', 88),
    ('渡辺あゆみ', 'watanabe@example.com', '2024-01-20 13:20:00', 'active', 94),
    ('伊藤雄一', 'ito@example.com', '2024-01-21 15:10:00', 'inactive', 72),
    ('中村真理', 'nakamura@example.com', '2024-01-22 08:45:00', 'active', 81),
    ('小林大輔', 'kobayashi@example.com', '2024-01-23 12:15:00', 'active', 89),
    ('加藤さくら', 'kato@example.com', '2024-01-24 17:30:00', 'pending', 76)
ON CONFLICT (email) DO NOTHING;