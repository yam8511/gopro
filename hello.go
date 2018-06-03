package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	parentCtx := context.Background()
	var key interface{} = "z"
	cancelCtx, cancel := context.WithCancel(parentCtx)
	valueCtx := context.WithValue(cancelCtx, key, "OK")
	go func() {
		// Pass a context with a timeout to tell a blocking function that it
		// should abandon its work after the timeout elapses.
		ctx, cancel := context.WithTimeout(valueCtx, 5000*time.Millisecond)
		defer cancel()

		select {
		case <-time.After(1 * time.Second):
			fmt.Println("overslept", ctx.Value(key))
		case <-ctx.Done():
			fmt.Println(ctx.Err(), ctx.Value(key)) // prints "context deadline exceeded"
		}
	}()

	go func() {
		d := time.Now().Add(1 * time.Second)
		ctx, cancel := context.WithDeadline(valueCtx, d)
		defer cancel()

		select {
		case <-time.After(2 * time.Second):
			fmt.Println("oversleep", ctx.Value(key))
		case <-ctx.Done():
			fmt.Println(ctx.Err(), ctx.Value(key))
		}
	}()

	time.Sleep(time.Second * 1)
	cancel()
	time.Sleep(time.Second * 1)
}
