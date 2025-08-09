## 実行方法

### 初回セットアップ
```bash
# イメージのビルド
docker-compose build

# WordCloud生成
docker-compose up

# 上記をターミナルで実行する

# 一度実行したビルドの各ステップ（キャッシュ）を無視してビルドし直す
docker-compose build --no-cache

# ワードクラウドの色を変える（サンプル参考）
# https://matplotlib.org/stable/users/explain/colors/colormaps.html