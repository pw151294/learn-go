package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key", "value")
	ctx, cancel := context.WithTimeout(ctx, 2*time.Hour)
	fmt.Println(ctx)
	fmt.Println(ctx.Value("key"))
	defer cancel()
}
