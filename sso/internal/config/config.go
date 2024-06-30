package config

import "time"

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTl time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
