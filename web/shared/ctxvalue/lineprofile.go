package ctxvalue

import (
	"context"
	"errors"

	"github.com/maguro-alternative/remake_bot/pkg/line"
)

const lineProfileKey contextKey = "lineProfile"

// ContextWithLineProfile sets the LineProfile in the context.
func ContextWithLineProfile(parent context.Context, profile *line.LineProfile) context.Context {
	return context.WithValue(parent, lineProfileKey, profile)
}

// LineProfileFromContext extracts the LineProfile from the context.
func LineProfileFromContext(ctx context.Context) (*line.LineProfile, error) {
	v := ctx.Value(lineProfileKey)
	profile, ok := v.(*line.LineProfile)
	if !ok {
		return nil, errors.New("profile not found")
	}
	return profile, nil
}
