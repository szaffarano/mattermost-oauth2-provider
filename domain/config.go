package domain

import "time"

// Config contains the server configuration data
type Config struct {
	Client struct {
		ID     string
		Secret string
		Domain string
	}
	JWTKey                  []byte
	GenerateRefresh         bool
	AuthorizeCodeExpiration time.Duration
	AccessTokenExpiration   time.Duration
	RefreshTokenExpiration  time.Duration

	MattermostURL string
	SUAURL        string
}
