package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	SSHPort       int    `mapstructure:"ssh_port"`
	ListenAddress string `mapstructure:"listen_address"`
	HostPrivKey   string `mapstructure:"host_priv_key"`
	AuthnURL      string `mapstructure:"authn_url"`
	LogLevel      string `mapstructure:"log_level"`
	LogFormat     string `mapstructure:"log_format"`
}

func Load(configFile string) (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("ssh_port", 2222)
	v.SetDefault("listen_address", "0.0.0.0")
	v.SetDefault("host_priv_key", "/etc/ssh/ssh_host_rsa_key")
	v.SetDefault("authn_url", "http://localhost:8080/authenticate")
	v.SetDefault("log_level", "info")
	v.SetDefault("log_format", "text")

	// Load config file if specified
	if configFile != "" {
		v.SetConfigFile(configFile)
		v.SetConfigType("yaml")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	// Read from environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("SSHSERVER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
