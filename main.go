package main

import (
	"context"
	"fmt"

	"github.com/go/stripe-gateway/cmd/api/dependencies"
)

func main() {
	deps, err := dependencies.New()
	if err != nil {
		handlePanic(err, "error setting up dependencies")
	}

	if err := deps.DB.Init(context.Background()); err != nil {
		handlePanic(err, "error initializing the DB")
	}

	if err := deps.API.Start(); err != nil {
		handlePanic(err, "error starting the API")
	}
}

func handlePanic(err error, errorMsg string) {
	panic(fmt.Errorf("%s: %s", errorMsg, err))
}
