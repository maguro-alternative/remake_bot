FROM golang:1.23.0-bookworm AS voicevox_setup

# Allow overriding the voicevox core asset URL at build time:
ARG VOICEVOX_VERSION="0.14.1"
ENV VOICEVOX_VERSION=${VOICEVOX_VERSION}

ARG TARGETARCH

WORKDIR /opt/voicevox

RUN apt-get -y update && apt-get install -y curl unzip file && rm -rf /var/lib/apt/lists/*

# Download voicevox_core (supports zip, tar.gz, or direct binary) with retries and debug output
# 1. ARM64(aarch64) 用の 0.14.1 バイナリをダウンロード
# URL 内の "x64" を "arm64" に変更するのが肝です
RUN if [ "$TARGETARCH" = "amd64" ]; then ARCH="x64"; else ARCH="arm64"; fi && \
    && wget https://github.com/VOICEVOX/voicevox_core/releases/download/0.14.1/voicevox_core-linux-${ARCH}-cpu-0.14.1.zip \
    && unzip voicevox_core-linux-${ARCH}-cpu-0.14.1.zip \
    && mv voicevox_core-linux-${ARCH}-cpu-0.14.1 core_files
RUN set -eux; \
    cp -a /voicevox_core/core_files/include/. /usr/local/include/ || true; \
    cp -a /voicevox_core/core_files/lib/. /usr/local/lib/ || true; \
    ldconfig || true

# ============================================================
# Stage 2: Go Builder with CGO support
# ============================================================
FROM golang:1.23.0-bookworm AS builder

# Install CGO dependencies and ffmpeg
RUN apt-get -y update && apt-get -y install locales && apt-get -y upgrade && \
    apt-get install -y ffmpeg build-essential pkg-config libopus-dev && \
    localedef -f UTF-8 -i ja_JP ja_JP.UTF-8 && \
    rm -rf /var/lib/apt/lists/*

ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8
ENV TZ JST-9
ENV TERM xterm

# Copy voicevox_core from setup stage
COPY --from=voicevox_setup /opt/voicevox /voicevox_core

# Set up library environment for CGO
ENV LD_LIBRARY_PATH=/voicevox_core:$LD_LIBRARY_PATH
ENV CGO_CFLAGS="-I/usr/local/include -I/voicevox_core/core_files"
ENV CGO_LDFLAGS="-L/usr/local/lib -L/voicevox_core/core_files -Wl,-rpath,/usr/local/lib"
ENV CGO_ENABLED=1

# Copy source code
COPY go.mod go.sum /root/src/
COPY pkg/ /root/src/pkg/
COPY core/ /root/src/core/
COPY bot/ /root/src/bot/
COPY web/ /root/src/web/
COPY repository/ /root/src/repository/
COPY tasks/ /root/src/tasks/
COPY testutil/ /root/src/testutil/

WORKDIR /root/src

# Configure Go for private repositories  
ENV GOPRIVATE=github.com/maguro-alternative/line-works-sdk-go
ENV GOPROXY=direct
ENV GOSUMDB=off

# Configure git for GitHub access (requires GITHUB_TOKEN build arg)
ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

# Download dependencies
RUN go mod download

# Build the application with CGO
RUN go build -o ./main ./core/main.go

# ============================================================
# Stage 3: Runtime image
# ============================================================
FROM ubuntu:24.04

RUN apt-get -y update && apt-get -y install locales && apt-get install -y \
    ffmpeg \
    curl \
    ca-certificates \
    && localedef -f UTF-8 -i ja_JP ja_JP.UTF-8 \
    && rm -rf /var/lib/apt/lists/*

ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8
ENV TZ JST-9
ENV TERM xterm

# Copy Go binary from builder stage
COPY --from=builder /root/src/main /app/main

# Copy voicevox_core library from setup stage
COPY --from=voicevox_setup /opt/voicevox /voicevox_core

# Set up library path
ENV LD_LIBRARY_PATH=/voicevox_core:$LD_LIBRARY_PATH

# Create startup script
RUN echo '#!/bin/bash\n\
set -e\n\
\n\
# Use Railway PORT environment variable or default to 5000\n\
PORT=${PORT:-5000}\n\
\n\
# Start main application with PORT environment variable\n\
echo "Starting application on port $PORT..."\n\
export PORT=$PORT\n\
exec /app/main\n\
' > /app/entrypoint.sh && chmod +x /app/entrypoint.sh

WORKDIR /app

EXPOSE 5000
CMD ["/app/entrypoint.sh"]
