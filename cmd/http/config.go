package main

import (
	"github.com/rendyananta/example-online-book-store/internal/config"
)

type BinaryConfig struct {
	HTTP HTTPConfig
	App  config.App
}

type HTTPConfig struct {
	ListenPort int
}

func loadConfig() BinaryConfig {
	return BinaryConfig{
		HTTP: HTTPConfig{
			ListenPort: config.LoadFromEnvInt("HTTP_LISTEN_PORT", defaultHTTPListenPort),
		},
		App: config.LoadAppConfig(),
	}
}
