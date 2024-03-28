package ctxvalue

import (
	"context"
	"errors"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

const lineUserKey contextKey = "lineUser"

// ContextWithLineUser sets the LineUser in the context.
func ContextWithLineUser(parent context.Context, user *model.LineOAuthSession) context.Context {
	return context.WithValue(parent, lineUserKey, user)
}

// LineUserFromContext extracts the LineUser from the context.
func LineUserFromContext(ctx context.Context) (*model.LineOAuthSession, error) {
	v := ctx.Value(lineUserKey)
	user, ok := v.(*model.LineOAuthSession)
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
