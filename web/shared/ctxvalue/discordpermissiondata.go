package ctxvalue

import (
	"context"
	"errors"

	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

const permissionKey contextKey = "discordPermissionData"

// ユーザー情報をコンテキストにセット
func ContextWithDiscordPermission(parent context.Context, user *model.DiscordPermissionData) context.Context {
    return context.WithValue(parent, permissionKey, user)
}

// ユーザー情報をコンテキストから取り出す
func DiscordPermissionFromContext(ctx context.Context) (*model.DiscordPermissionData, error) {
    v := ctx.Value(permissionKey)
    user, ok := v.(*model.DiscordPermissionData)
    if !ok {
        return nil, errors.New("permission not found")
    }
    return user, nil
}