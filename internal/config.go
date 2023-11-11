package internal

import "fmt"

type Config struct {
	RelayHost  string `env:"SERVER_HOST" env-default:"localhost"`
	RelayPort  string `env:"SERVER_PORT" env-default:"25"`
	Host       string `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
	Port       string `env:"SMTP_PORT" env-default:"587"`
	User       string `env:"SMTP_USER"`
	Pass       string `env:"SMTP_PASS"`
	StartTLS   bool   `env:"SMTP_TLS" env-default:"true"`
	Auth       bool   `env:"SMTP_AUTH" env-default:"true"`
	SlackToken string `env:"SLACK_TOKEN" env-required:"true"`
}

func (c *Config) Insecure() bool {
	return !c.StartTLS
}

func (c *Config) RelayAddr() string {
	return fmt.Sprintf("%s:%s", c.RelayHost, c.RelayPort)
}

func (c *Config) SMTPAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
