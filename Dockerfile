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
RUN if [ -n "$GITHUB_TOKEN" ]; then \
        git config --global url."https://x-access-token:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; \
        echo "Git configured with GitHub token"; \
    else \
        echo "Warning: No GITHUB_TOKEN provided. This may cause build failures for private repositories."; \
        echo "Use: docker build --build-arg GITHUB_TOKEN=your_token_here ."; \
    fi

# Download dependencies (will fetch from GitHub directly)
RUN go mod download

# Build the application
RUN go build -o ./main ./core/main.go
