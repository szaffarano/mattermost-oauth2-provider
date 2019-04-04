package domain

import "time"

// Config contains the server configuration data
type Config struct {
	Client struct {
		ID     string
		Secret string
		Domain string
	}

	Oauth struct {
		JWTKey                  string
		GenerateRefresh         bool
		AuthorizeCodeExpiration time.Duration
		AccessTokenExpiration   time.Duration
		RefreshTokenExpiration  time.Duration
	}

	App struct {
		Port          int32
		MattermostURL string
		SUAURL        string
		PublicKey     string
		Verbose       bool `yaml:"-"`
	}
}

// Validate validates the configuration
func (cfg Config) Validate() error {
	return nil
}
