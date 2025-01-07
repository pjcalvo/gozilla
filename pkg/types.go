package pkg

import (
	"context"
)

type ExecuteFunc func(context.Context) (any, error)
type ExpectFunc func(context.Context, any, error) (error)