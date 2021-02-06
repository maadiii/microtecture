package config

import (
	"os"

	"github.com/pkg/errors"
)

type refreshToken struct {
	Secret    string `yaml:"secret"`
	Algorithm string `yaml:"algorithm"`
	MaxAge    uint   `yaml:"max_age"`
	Secure    bool   `yaml:"secure"`
	HTTPOnly  bool   `yaml:"httponly"`
	Path      string `yaml:"path"`
}

type jwt struct {
	Secret       string       `yaml:"secret"`
	Algorithm    string       `yaml:"algorithm"`
	MaxAge       uint         `yaml:"max_age"`
	HTTPOnly     bool         `yaml:"httponly"`
	RefreshToken refreshToken `yaml:"refresh_token"`
}

type ApplicationConfig struct {
	IsDevelopment bool
	SecretKey     string `yaml:"secret_key"`
	JWT           jwt    `yaml:"jwt"`
	Port          int16  `yaml:"port"`
}

func (self *ApplicationConfig) Init() error {
	if os.Getenv(ENVIRONMENT_NAME) == "dev" {
		self.IsDevelopment = true
	}

	if err := marshal(self); err != nil {
		return err
	}

	if len(self.SecretKey) < 32 {
		return errors.New("secret_key is not set in config file or lesser than 32.")
	}

	if len(self.JWT.Secret) < 8 {
		return errors.New("jwt.secret is not set in config file or lesser than 8.")
	}

	if self.JWT.Algorithm != HS256 &&
		self.JWT.Algorithm != HS384 &&
		self.JWT.Algorithm != HS512 {
		return errors.New("jwt.algorithm is not set in config file or not in (HS256, HS384, HS512).")
	}

	if self.JWT.MaxAge == 0 {
		return errors.New("jwt.max_age is not set in config file.")
	}

	if len(self.JWT.RefreshToken.Secret) < 8 {
		return errors.New("jwt.refresh_token.secret is not set in config file or lesser than 8.")

	}

	if self.JWT.RefreshToken.Algorithm != HS256 &&
		self.JWT.RefreshToken.Algorithm != HS384 &&
		self.JWT.RefreshToken.Algorithm != HS512 {
		return errors.New(
			"jwt.refresh_token.algorithm is not set in config file or not in (HS256, HS384, HS512).",
		)
	}

	if self.JWT.RefreshToken.MaxAge == 0 {
		return errors.New("jwt.refresh_token.max_age is not set in config file.")
	}

	if self.JWT.RefreshToken.Path == "" {
		return errors.New("jwt.refresh_token.path is not set in config file.")
	}

	if self.Port == 0 {
		return errors.New("http_port is not set in config file.")
	}

	return nil
}
