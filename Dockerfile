FROM golang:1.23.0-bullseye AS builder

RUN apt-get -y update && apt-get -y install locales && apt-get -y upgrade && apt-get install -y ffmpeg &&\
    localedef -f UTF-8 -i ja_JP ja_JP.UTF-8
ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8
ENV TZ JST-9
ENV TERM xterm

# Copy source code (no need for vendor directory anymore)
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

# Build the application
RUN go build -o ./main ./core/main.go

# ============================================================
# Stage 2: voicevox_core installer
# ============================================================
FROM python:3.11-bookworm AS voicevox_installer

WORKDIR /opt/voicevox

# Install voicevox_core dependencies
RUN apt-get -y update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Download voicevox_core (latest CPU version - direct binary download)
# Note: The download URLs are for direct binaries, not zip files
RUN curl -L "https://github.com/VOICEVOX/voicevox_core/releases/download/0.16.3/download-linux-x64" -o download && \
    chmod +x download && \
    mkdir models && \
    # Download default models
    curl -L -o models/0.vvm https://raw.githubusercontent.com/VOICEVOX/voicevox_vvm/main/vvms/0.vvm && \
    # Run with automatic acceptance of terms (non-interactive)
    echo "y" | ./download --exclude models || true

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

# Copy voicevox_core from installer stage
COPY --from=voicevox_installer /opt/voicevox/voicevox_core /app/voicevox_core

# Create startup script
RUN echo '#!/bin/bash\n\
set -e\n\
\n\
# Use Railway PORT environment variable or default to 5000\n\
PORT=${PORT:-5000}\n\
\n\
# Check if voicevox should be disabled\n\
ENABLE_VOICEVOX=${ENABLE_VOICEVOX:-true}\n\
\n\
if [ "$ENABLE_VOICEVOX" = "true" ]; then\n\
  # Start voicevox_core in background\n\
  echo "Starting VOICEVOX Core..."\n\
  /app/voicevox_core --port 50021 --cpu_num_threads 2 > /tmp/voicevox.log 2>&1 &\n\
  VOICEVOX_PID=$!\n\
  \n\
  # Wait for voicevox_core to be ready\n\
  echo "Waiting for VOICEVOX Core to start..."\n\
  for i in {1..30}; do\n\
    if curl -s http://localhost:50021/version > /dev/null 2>&1; then\n\
      echo "VOICEVOX Core started successfully"\n\
      break\n\
    fi\n\
    if [ $i -eq 30 ]; then\n\
      echo "Warning: Failed to start VOICEVOX Core, continuing without it"\n\
      cat /tmp/voicevox.log\n\
      break\n\
    fi\n\
    sleep 1\n\
  done\n\
else\n\
  echo "VOICEVOX Core disabled"\n\
fi\n\
\n\
# Start main application with PORT environment variable\n\
echo "Starting application on port $PORT..."\n\
export PORT=$PORT\n\
exec /app/main\n\
' > /app/entrypoint.sh && chmod +x /app/entrypoint.sh

WORKDIR /app

EXPOSE 5000
CMD ["/app/entrypoint.sh"]
