package main

import (
	"context"
	"log"
	"no-q-solution/bootstrap"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	bootstrap.Start(ctx)
}
