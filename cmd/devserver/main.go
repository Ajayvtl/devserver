package main

import (
	"context"
	"os"

	"github.com/Ajayvtl/devserver/internal/logger"
)

func main() {
	log := logger.New(logger.Config{Service: "devserver"})

	if err := Run(context.Background(), os.Args[1:]); err != nil {
		log.Fatal().Err(err).Msg("devserver exited")
	}
}
