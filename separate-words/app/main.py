# MeCab + pandas + wordcloudの処理

import MeCab
import pandas as pd
from wordcloud import WordCloud
import matplotlib.pyplot as plt
import japanize_matplotlib
import os

# MeCab Taggerの初期化
mecab = MeCab.Tagger()

# テスト用の文章リスト
texts = [
    "Peach.Techは成蹊大学の公認サークルです。大学設立当初、Peach.Techは存在していませんでした。",
    "~~あとで編集し直して公開に変える~~　🆗",
    "山月記ってみんな知ってる話じゃなかったんだ",
    "iphoneとandroidで絵文字のタッチ全然違うから、可愛いつもりで使ってて全然可愛くないことが頻発する",
    "さん、Peach.Techへようこそ！早速ですが以下のチャンネルで自己紹介をお願いします！🍑 https://discord.com/channels/",
]

# 結果格納用
data = []

for sentence in texts:
    words, roots, parts = [], [], []
    node = mecab.parseToNode(sentence) # nodeは文節のこと
    while node:
        surface = node.surface # 表層形
        features = node.feature.split(",") # mecabの出力結果をコンマ区切りで取得
        base = features[6] if len(features) > 6 else "*" # 原形
        pos = features[0] # 品詞
        if surface:
            words.append(surface)
            roots.append(base)
            parts.append(pos)
        node = node.next
    data.append({"sentence": sentence, "words": words, "root": roots, "part": parts})

# 解析結果をDataFrameに変換
df = pd.DataFrame(data)

# 自動詞を新リスト（filtered_words）に格納
filtered_words = []
for i, row in df.iterrows():
    for root, part in zip(row["root"], row["part"]):
        if part not in ["助詞", "助動詞"]:
            filtered_words.append(root)

# intransitiveの要素をスペース区切りで連結
text_for_wc = " ".join(filtered_words)

# 出力ディレクトリを作成
from datetime import datetime
OUTPUT_DIR = os.getenv("OUTPUT_DIR", "/app/output")
os.makedirs(OUTPUT_DIR, exist_ok=True)

import numpy as np
from PIL import Image

mask_image = np.array(Image.open("/app/logo/PeachTech_black.png"))

# ワードクラウド生成（フォント指定なし）
wordcloud = WordCloud(
    font_path="/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
    background_color="white",
    mask=mask_image,
    colormap="Paired",
    width=800,
    height=800
).generate(text_for_wc)

# タイムスタンプ付きのファイル名を生成
timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
output_filename = f"wordcloud_output_{timestamp}.png"
output_path = os.path.join(OUTPUT_DIR, output_filename)
wordcloud.to_file(output_path)
print(f"✅ WordCloud画像を保存しました → {output_path}")