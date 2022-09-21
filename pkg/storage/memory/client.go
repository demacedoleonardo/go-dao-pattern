package memory

import (
	"context"
)

type (
	Memory map[string]interface{}

	Client interface {
		Get(ctx *context.Context, key string) (interface{}, error)
		Save(ctx *context.Context, key string, value interface{}) error
	}
)
