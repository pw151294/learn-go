package main

import (
	"context"
	"fmt"
)

type contextKey string

func main() {
	f := func(ctx context.Context, k contextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Printf("found value for key %v: %v\n", k, v)
		} else {
			fmt.Printf("not found value: %v\n", k)
		}
	}

	ctx := context.WithValue(context.Background(), contextKey("language"), "Go")
	f(ctx, "color")
	f(ctx, "language")
}
