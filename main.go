package main

import (
	"context"
	"log/slog"
	"os"
)

const DEFAULT_CONFIG_LOCATION = "/etc/temps_de_funcionament/config.yaml"

func main() {
	config, err := GetConfiguration(getConfigLocation(os.Args))
	if err != nil {
		slog.With("error", err).Error("Error getting configuration")
		return
	}
	slackNotifier := NewSlackNotifier(os.Getenv("TOKEN"))
	t := NewTempsDeFuncionament(config, slackNotifier)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Start(ctx)
	<-ctx.Done()
}

func getConfigLocation(args []string) string {
	if len(args) == 1 {
		return DEFAULT_CONFIG_LOCATION
	}
	configLocation := os.Args[1]
	if configLocation == "" {
		configLocation = DEFAULT_CONFIG_LOCATION
	}
	return configLocation
}
