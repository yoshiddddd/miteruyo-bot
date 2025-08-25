## 1. Adminerを使用（推奨・簡単）:
コンテナを起動
  docker-compose up -d

  ブラウザで http://localhost:8080 
  にアクセス
    ログイン情報:
    システム: PostgreSQL
    サーバ: postgres
    ユーザ名: admin
    パスワード: password123
    データベース: datacollect

  ## 2. psqlコマンドを使用:
  PostgreSQLコンテナに接続
  docker exec -it data-collect-postgres psql -U admin -d datacollect

 SQLクエリを実行
  SELECT * FROM users;

  ## 3. 外部のPostgreSQLクライアントを使用:
  - DBeaver、pgAdmin、TablePlusなどのGUIクライ
  アント
  - 接続情報: localhost:5432, ユーザー: admin,
   パスワード: password123