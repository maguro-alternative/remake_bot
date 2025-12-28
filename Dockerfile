FROM golang:1.23.0-bullseye AS builder

RUN apt-get -y update && apt-get -y install locales && apt-get -y upgrade && apt-get install -y ffmpeg &&\
    localedef -f UTF-8 -i ja_JP ja_JP.UTF-8
ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8
ENV TZ JST-9
ENV TERM xterm

# ./root/src ディレクトリを作成 ホームのファイルをコピーして、移動
RUN mkdir -p /root/src
COPY . /root/src
WORKDIR /root/src

# Initialize git submodules if they exist
RUN if [ -f .gitmodules ]; then \
    git config --global --add safe.directory /root/src && \
    git config --global --add safe.directory /root/src/vendor/line-works-sdk-go && \
    git submodule update --init --recursive; \
    fi

# Docker内で扱うffmpegをインストール
RUN go mod download && \
    go build -o ./main ./core/main.go
