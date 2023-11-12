package internal

import "fmt"

type Config struct {

	// ServerHost and ServerPort are used for listening incoming connections.
	// By default, it allows the server to respond to requests from any
	// available network interface.
	ServerHost string `env:"SERVER_HOST" env-default:"0.0.0.0"`
	ServerPort string `env:"SERVER_PORT" env-default:"25"`
	Host       string `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
	Port       string `env:"SMTP_PORT" env-default:"587"`
	User       string `env:"SMTP_USER"`
	Pass       string `env:"SMTP_PASS"`
	StartTLS   bool   `env:"SMTP_TLS" env-default:"true"`
	Auth       bool   `env:"SMTP_AUTH" env-default:"true"`

	// Slack token for sending notifications.
	// https://api.slack.com/authentication/token-types#granular_bot
	SlackToken string `env:"SLACK_TOKEN"`
}

func (c *Config) Insecure() bool {
	return !c.StartTLS
}

func (c *Config) RelayAddr() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

func (c *Config) SMTPAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
