package main

import (
	"context"
)

type Workflow interface {
	Name() string
	Run(ctx context.Context) error
}
