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

# Setup Git configuration for submodules (needed for private repos)
RUN git config --global user.email "docker@example.com" && \
    git config --global user.name "Docker Build"

# Initialize git repository and restore submodules
# This step requires the .git directory to be copied
RUN if [ -d .git ]; then \
        git config --global --add safe.directory /root/src && \
        git config --global --add safe.directory /root/src/vendor/line-works-sdk-go && \
        git submodule update --init --recursive; \
    else \
        echo "Warning: .git directory not found. Submodules will not be initialized."; \
        echo "Make sure vendor/line-works-sdk-go is included in the Docker build context."; \
    fi

# Verify that the vendor directory exists
RUN if [ ! -d vendor/line-works-sdk-go/pkg/lineworks ]; then \
        echo "Error: LINE Works SDK not found at vendor/line-works-sdk-go/pkg/lineworks"; \
        echo "Please ensure the submodule is properly initialized or the vendor directory is included."; \
        exit 1; \
    fi

# Docker内で扱うffmpegをインストール
RUN go mod download && \
    go build -o ./main ./core/main.go
