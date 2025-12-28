FROM golang:1.23.0-bullseye AS builder

RUN apt-get -y update && apt-get -y install locales && apt-get -y upgrade && apt-get install -y ffmpeg &&\
    localedef -f UTF-8 -i ja_JP ja_JP.UTF-8
ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8
ENV TZ JST-9
ENV TERM xterm

# Step 1: Copy everything except .git to reduce build context
COPY go.mod go.sum /root/src/
COPY pkg/ /root/src/pkg/
COPY core/ /root/src/core/
COPY bot/ /root/src/bot/
COPY web/ /root/src/web/
COPY repository/ /root/src/repository/
COPY tasks/ /root/src/tasks/
COPY testutil/ /root/src/testutil/

# Step 2: Explicitly copy vendor directory to ensure it's included
COPY vendor/ /root/src/vendor/

WORKDIR /root/src

# Debug: Check if vendor directory is properly copied
RUN echo "=== Checking vendor directory after COPY ===" && \
    echo "Vendor directory exists:" && \
    ls -la vendor/ && \
    echo "line-works-sdk-go directory exists:" && \
    ls -la vendor/line-works-sdk-go/ && \
    echo "Total files in line-works-sdk-go:" && \
    find vendor/line-works-sdk-go -type f | wc -l && \
    echo "Go files in line-works-sdk-go:" && \
    find vendor/line-works-sdk-go -name "*.go" | head -5 || echo "No Go files found"

# Verify that the vendor directory exists
RUN if [ ! -d vendor/line-works-sdk-go/pkg/lineworks ]; then \
        echo "Error: LINE Works SDK not found at vendor/line-works-sdk-go/pkg/lineworks"; \
        echo "Available vendor contents:"; \
        ls -la vendor/ || echo "No vendor directory"; \
        if [ -d vendor/line-works-sdk-go ]; then \
            echo "line-works-sdk-go contents:"; \
            ls -la vendor/line-works-sdk-go/; \
        fi; \
        echo ""; \
        echo "To fix this issue, run locally:"; \
        echo "  git submodule update --init --recursive"; \
        echo "  docker build ."; \
        exit 1; \
    else \
        echo "LINE Works SDK found successfully"; \
        echo "SDK contents:"; \
        ls -la vendor/line-works-sdk-go/pkg/lineworks/; \
    fi

# Docker内で扱うffmpegをインストール
RUN go mod download && \
    go build -mod=mod -o ./main ./core/main.go
