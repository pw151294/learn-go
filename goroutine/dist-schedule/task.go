package main

import (
	"context"
)

type Task interface {
	ID() string
	Execute(ctx context.Context) (*AdResult, error)
}

type AdResult struct {
	Result interface{}
}

type Error struct {
	taskID string
	err    error
}

type AdTask struct {
	id     string
	adFunc func(ctx context.Context) (*AdResult, error)
}

func (t *AdTask) ID() string {
	return t.id
}

func (t *AdTask) Execute(ctx context.Context) (*AdResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return t.adFunc(ctx)
	}
}
