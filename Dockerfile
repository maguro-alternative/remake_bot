FROM golang:1.23.0-bookworm AS voicevox_setup

WORKDIR /opt/voicevox

RUN apt-get -y update && apt-get install -y curl unzip && rm -rf /var/lib/apt/lists/*

# Download voicevox_core library (Linux x86_64)
RUN curl -L "https://github.com/VOICEVOX/voicevox_core/releases/download/0.16.3/voicevox_core-0.16.3-linux-x64.zip" -o voicevox.zip && \
    unzip voicevox.zip && \
    rm voicevox.zip && \
    ls -la

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
