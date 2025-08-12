## 実行方法

### 初回セットアップ
```bash
# イメージのビルド
docker-compose build
## キャッシュを無視してビルドし直す場合（コードを更新したときなど）
docker-compose build --no-cache

# コンテナを起動
docker-compose up

# 上記をターミナルで実行する


# ワードクラウドの色を変える（サンプル参考）
# https://matplotlib.org/stable/users/explain/colors/colormaps.html