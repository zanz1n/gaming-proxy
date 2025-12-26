package main

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

var valueCfg *Config = nil

func GetConfig() *Config {
	return valueCfg
}

func InitConfig() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := new(Config)
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return err
	}

	valueCfg = cfg
	return nil
}

type Config struct {
	Proxied    ProxiedConfig    `env:", prefix=PROXIED_"`
	Cloudflare CloudflareConfig `env:", prefix=CLOUDFLARE_"`
	Ngrok      NgrokConfig      `env:", prefix=NGROK_"`
}

type ProxiedConfig struct {
	Host string `env:"HOST, default=127.0.0.1"`
	Port uint16 `env:"PORT, required"`
}

type CloudflareConfig struct {
	Token     string `env:"TOKEN, required"`
	ZoneID    string `env:"ZONE_ID, required"`
	Subdomain string `env:"SUB_DOMAIN, required"`
	Service   string `env:"SERVICE, default=minecraft"`
	Protocol  string `env:"PROTOCOL, default=tcp"`
	Overwrite bool   `env:"OVERWRITE, default=true"`
}

type NgrokConfig struct {
	Token string `env:"TOKEN, required"`
}
