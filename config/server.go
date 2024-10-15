package config

import "time"

type ServerConfig struct {
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	WriteTimeout     uint   `mapstructure:"write-timeout"`
	ReadTimeout      uint   `mapstructure:"read-timeout"`
	IdleTimeout      uint   `mapstructure:"idle-timeout"`
	MaxContentLength uint32 `mapstructure:"max-content-length"`
	AccessToken      string `mapstructure:"access-token"`
}

type ParsedServerConfig struct {
	Host             string
	Port             int
	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	IdleTimeout      time.Duration
	MaxContentLength uint32
	AccessToken      string
}

func (c *ServerConfig) Parse() (*ParsedServerConfig, error) {
	// TODO Add some validations
	return &ParsedServerConfig{
		Host:             c.Host,
		Port:             c.Port,
		WriteTimeout:     time.Duration(c.WriteTimeout) * time.Second,
		ReadTimeout:      time.Duration(c.ReadTimeout) * time.Second,
		IdleTimeout:      time.Duration(c.IdleTimeout) * time.Second,
		MaxContentLength: c.MaxContentLength,
		AccessToken:      c.AccessToken,
	}, nil
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:             "127.0.0.1",
		Port:             9791,
		WriteTimeout:     15,
		ReadTimeout:      15,
		IdleTimeout:      120,
		MaxContentLength: 8192,
		AccessToken:      "",
	}
}
