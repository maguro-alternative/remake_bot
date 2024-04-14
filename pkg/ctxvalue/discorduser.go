package ctxvalue

import (
	"context"
	"errors"

	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

type contextKey string

const discordUserKey contextKey = "discordUser"

// ユーザー情報をコンテキストにセット
func ContextWithDiscordUser(parent context.Context, user *model.DiscordOAuthSession) context.Context {
	return context.WithValue(parent, discordUserKey, user)
}

// ユーザー情報をコンテキストから取り出す
func DiscordUserFromContext(ctx context.Context) (*model.DiscordOAuthSession, error) {
	v := ctx.Value(discordUserKey)
	user, ok := v.(*model.DiscordOAuthSession)
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
