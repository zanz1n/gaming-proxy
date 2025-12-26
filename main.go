package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

	if err := InitConfig(); err != nil {
		fatal("Config", err)
	}
}

func main() {
	s1 := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	mainCtx, mainCancel := context.WithCancel(context.Background())

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		sig := <-ch
		log.Printf("%s: shutting down ...\n", strings.ToUpper(sig.String()))
		cancel()
		mainCancel()
	}()

	start := time.Now()
	res, err := RunNgrok(ctx)
	if err != nil {
		fatal("Ngrok", err)
	}
	defer res.Close()

	slog.Info("Ngrok is setup",
		"addr", fmt.Sprintf("%s:%d", res.Host, res.Port),
		"took", time.Since(start).Round(time.Microsecond),
	)

	start = time.Now()
	err = RunCloudflare(ctx, res.Host, res.Port)
	if err != nil {
		fatal("Cloudflare", err)
	}

	slog.Info("Cloudflare is setup",
		"took", time.Since(start).Round(time.Millisecond),
	)

	slog.Info("Running",
		"took",
		time.Since(s1).Round(time.Millisecond),
	)

	<-mainCtx.Done()
}

func fatal(name string, err error) {
	fmt.Printf("\x1b[1m\x1b[31m%s error:\x1b[0m %s\n", name, err.Error())
	os.Exit(1)
}
