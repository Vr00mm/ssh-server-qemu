package main

import (
	"fmt"
	"os"

	"github.com/phuslu/log"
	cli "github.com/urfave/cli/v2"

	"yourusername/sshserver/internal/config"
	"yourusername/sshserver/internal/logging"
	"yourusername/sshserver/internal/server"
)

func main() {
	app := &cli.App{
		Name:  "sshserver",
		Usage: "A configurable SSH server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "SSH server port",
				EnvVars: []string{"SSHSERVER_SSH_PORT"},
			},
			&cli.StringFlag{
				Name:    "listen",
				Aliases: []string{"l"},
				Usage:   "Listen address",
				EnvVars: []string{"SSHSERVER_LISTEN_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "privkey",
				Aliases: []string{"k"},
				Usage:   "Host private key file",
				EnvVars: []string{"SSHSERVER_HOST_PRIV_KEY"},
			},
			&cli.StringFlag{
				Name:    "authn-url",
				Aliases: []string{"a"},
				Usage:   "Authentication server URL",
				EnvVars: []string{"SSHSERVER_AUTHN_URL"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "Log level (trace, debug, info, warn, error, fatal)",
				EnvVars: []string{"SSHSERVER_LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-format",
				Usage:   "Log format (json or text)",
				EnvVars: []string{"SSHSERVER_LOG_FORMAT"},
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run application")
	}
}

func runServer(c *cli.Context) error {
	cfg, err := config.Load(c.String("config"))
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override config with command line flags
	if c.IsSet("port") {
		cfg.SSHPort = c.Int("port")
	}
	if c.IsSet("listen") {
		cfg.ListenAddress = c.String("listen")
	}
	if c.IsSet("privkey") {
		cfg.HostPrivKey = c.String("privkey")
	}
	if c.IsSet("authn-url") {
		cfg.AuthnURL = c.String("authn-url")
	}
	if c.IsSet("log-level") {
		cfg.LogLevel = c.String("log-level")
	}
	if c.IsSet("log-format") {
		cfg.LogFormat = c.String("log-format")
	}

	logging.Setup(cfg.LogLevel, cfg.LogFormat)

	sshServer, err := server.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create SSH server: %w", err)
	}

	return sshServer.Start()
}
