#!/bin/sh
# schema/schema.sql と DB の差分を atlas schema apply で適用する(Railway の pre-deploy から実行)。
# 必要な環境変数:
#   DATABASE_URL  : 適用先DB
#   ATLAS_DEV_URL : 差分計算用の空DB(適用先と同じサーバー上の atlas_dev など)
set -eu
cd "$(dirname "$0")/.."
: "${DATABASE_URL:?DATABASE_URL is required}"
: "${ATLAS_DEV_URL:?ATLAS_DEV_URL is required}"
exec atlas schema apply --env prod --auto-approve
