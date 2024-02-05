package config

import (
	"fmt"
	"sync"

	"github.com/maguro-alternative/remake_bot/bot/config/internal"

	"github.com/caarlos0/env/v7"
	"github.com/cockroachdb/errors"
)

var (
	once sync.Once
	cfg  *internal.Config
)

func init() {
	once.Do(MustInit)
}

func MustInit() {
	cfg = &internal.Config{}
	if err := env.Parse(cfg); err != nil {
		xerr := errors.Wrap(err, "failed to env parse: ")
		fmt.Printf("panic: %+v", xerr)
		panic(xerr)
	}
}

func Token() string {
	return cfg.Token
}
