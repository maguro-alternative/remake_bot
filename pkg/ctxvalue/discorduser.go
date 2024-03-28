package ctxvalue

import (
	"context"
	"errors"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

type contextKey string

const userKey contextKey = "discordUser"

// ユーザー情報をコンテキストにセット
func ContextWithDiscordUser(parent context.Context, user *model.DiscordUser) context.Context {
    return context.WithValue(parent, userKey, user)
}

// ユーザー情報をコンテキストから取り出す
func DiscordUserFromContext(ctx context.Context) (*model.DiscordUser, error) {
    v := ctx.Value(userKey)
    user, ok := v.(*model.DiscordUser)
    if !ok {
        return nil, errors.New("user not found")
    }
    return user, nil
}