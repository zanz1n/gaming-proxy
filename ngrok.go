package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"golang.ngrok.com/ngrok/v2"
)

type NgrokResult struct {
	Host  string
	Port  uint16
	Close func() error
}

func RunNgrok(ctx context.Context) (NgrokResult, error) {
	cfg := GetConfig()

	agent, err := ngrok.NewAgent(
		ngrok.WithAuthtoken(cfg.Ngrok.Token),
		ngrok.WithLogger(slog.Default()),
	)
	if err != nil {
		return NgrokResult{}, err
	}
	if err := agent.Connect(ctx); err != nil {
		return NgrokResult{}, err
	}

	upstreamAddr := fmt.Sprintf("tcp://%s:%d", cfg.Proxied.Host, cfg.Proxied.Port)
	fowarder, err := agent.Forward(ctx,
		ngrok.WithUpstream(upstreamAddr),
		ngrok.WithURL("tcp://"),
	)
	if err != nil {
		return NgrokResult{}, err
	}

	host, portstr, _ := strings.Cut(fowarder.URL().Host, ":")
	port, err := strconv.Atoi(portstr)
	if err != nil {
		return NgrokResult{}, fmt.Errorf("invalid ngrok response: %w", err)
	}

	closefn := func() error {
		err1 := fowarder.Close()
		err2 := agent.Disconnect()

		if err1 != nil {
			if err2 != nil {
				return errors.Join(err1, err2)
			}
			return err1
		}
		return err2
	}

	return NgrokResult{
		Host:  host,
		Port:  uint16(port),
		Close: closefn,
	}, nil
}
