package appinit

import (
	"context"

	"github.com/joho/godotenv"
)

func InitDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		func(ctx context.Context) error { return initConfig(ctx) },
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func initConfig(_ context.Context) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	return nil
}
