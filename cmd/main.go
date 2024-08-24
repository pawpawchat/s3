package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	mycfg "github.com/pawpawchat/s3/config"
	"github.com/pawpawchat/s3/internal/app"
)

func main() {
	var cfg *mycfg.Config
	mycfg.ConfigureLogger(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-exit
		cancel()
	}()

	app.Run(ctx, cfg)
}
