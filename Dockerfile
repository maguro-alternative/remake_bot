FROM golang:1.23.0-bookworm AS voicevox_setup

# Allow overriding the voicevox core asset URL at build time:
# docker build --build-arg VOICEVOX_ASSET_URL=<url> -t remake_bot:v1 .
ARG VOICEVOX_ASSET_URL="https://github.com/VOICEVOX/voicevox_core/releases/download/0.16.3/download-linux-x64"
ARG VOICEVOX_VERSION="0.16.3"
ENV VOICEVOX_ASSET_URL=${VOICEVOX_ASSET_URL}
ENV VOICEVOX_VERSION=${VOICEVOX_VERSION}

WORKDIR /opt/voicevox

RUN apt-get -y update && apt-get install -y curl unzip file && rm -rf /var/lib/apt/lists/*

# Download voicevox_core (supports zip, tar.gz, or direct binary) with retries and debug output
RUN set -eux; \
        OUT=/opt/voicevox/voicevox.asset; \
        curl -sSL -D /tmp/curl_headers.txt --retry 3 --retry-delay 5 -o "$OUT" "$VOICEVOX_ASSET_URL" || true; \
        echo "curl headers:"; cat /tmp/curl_headers.txt || true; \
        echo "downloaded file info:"; ls -la "$OUT" || true; \
        if [ ! -s "$OUT" ] || [ $(stat -c%s "$OUT" || echo 0) -lt 1000 ]; then \
                echo "Downloaded file looks too small or empty; dumping first 4KiB:"; \
                head -c 4096 "$OUT" || true; \
                false; \
        fi; \
        mimetype=$(file -b --mime-type "$OUT" || true); \
        echo "detected mime-type: $mimetype"; \
        case "$mimetype" in \
            application/zip) \
                unzip "$OUT" -d /opt/voicevox ;; \
            application/gzip|application/x-gzip) \
                tar -xzf "$OUT" -C /opt/voicevox ;; \
            application/x-tar) \
                tar -xf "$OUT" -C /opt/voicevox ;; \
            application/octet-stream|binary|application/x-pie-executable|application/x-executable|application/x-elf) \
                echo "Treating as binary; moving to /opt/voicevox/"; \
                mkdir -p /opt/voicevox && mv "$OUT" /opt/voicevox/voicevox_core_binary && chmod +x /opt/voicevox/voicevox_core_binary; \
                echo "Attempting to fetch header files for VOICEVOX ${VOICEVOX_VERSION}"; \
                curl -fsSL -o /opt/voicevox/voicevox_core.h "https://raw.githubusercontent.com/VOICEVOX/voicevox_core/v${VOICEVOX_VERSION}/include/voicevox_core.h" || echo "warning: couldn't fetch voicevox_core.h"; \
                curl -fsSL -o /opt/voicevox/version.h "https://raw.githubusercontent.com/VOICEVOX/voicevox_core/v${VOICEVOX_VERSION}/include/version.h" || true ;; \
            *) \
                echo "Unknown asset mime-type: $mimetype"; false ;; \
        esac; \
        rm -f /tmp/curl_headers.txt || true; \
        ls -la /opt/voicevox

# ============================================================
# Stage 2: Go Builder with CGO support
# ============================================================
FROM golang:1.23.0-bookworm AS builder

# Install CGO dependencies and ffmpeg
RUN apt-get -y update && apt-get -y install locales && apt-get -y upgrade && \
    apt-get install -y ffmpeg build-essential pkg-config && \
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
ENV CGO_CFLAGS="-I/voicevox_core"
ENV CGO_LDFLAGS="-L/voicevox_core"
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
