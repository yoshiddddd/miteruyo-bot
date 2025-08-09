FROM python:3.11-slim
# PythonでMecabを使うためのライブラリ
RUN pip install --no-cache-dir mecab-python3

# タイムゾーン設定の対話を回避
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Tokyo

# MeCabとその依存関係をインストール
RUN apt-get update && apt-get install -y \
    mecab \
    mecab-ipadic-utf8 \
    libmecab-dev \
    git \
    curl \
    file \
    xz-utils \
    make \
    build-essential \
    sudo \
    fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/*

# NEologd辞書をインストール
RUN git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git /tmp/neologd \
    && /tmp/neologd/bin/install-mecab-ipadic-neologd -n -y \
     # デフォルト辞書を NEologd に向ける（mecabrc を上書き）
    && echo "dicdir = $(mecab-config --dicdir)/mecab-ipadic-neologd" > /etc/mecabrc \
    && cp /etc/mecabrc /usr/local/etc/mecabrc \
    && rm -rf /tmp/neologd

# Dockerイメージに必要なファイルをコピー
COPY requirements.txt .
RUN pip install -r requirements.txt
# 出力ディレクトリを作成し、権限を設定
RUN mkdir -p /app/output && chmod 777 /app/output
COPY . .

CMD ["python", "app/main.py"]