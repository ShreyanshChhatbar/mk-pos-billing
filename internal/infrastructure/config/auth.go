package config

import "os"

type AuthConfig struct {
	JWTSecret string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
