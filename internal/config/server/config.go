// Package config provides configuration handling for the server, allowing
// configurations to be set via environment variables or command-line flags.
package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/eschao/config"
	"github.com/eschao/config/env"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mrkovshik/memento/internal/validation"
)

const (
	defaultConfigFilePath = ""
	defaultAddress        = "localhost:3200"
	defaultDBAddress      = "host=localhost port=5432 user=yandex password=yandex dbname=final sslmode=disable"
	defaultCryptoKey      = "MySuperAwesomeCryptoKey"
	defaultTokenExpiry    = time.Hour * 24 * 7
)

var k = koanf.New(".")

// ServerConfig holds the configuration settings for the server.
type ServerConfig struct {
	Address             string `env:"ADDRESS" json:"address"`
	addressIsSet        bool
	DBAddress           string `env:"DATABASE_DSN" json:"database_dsn"`
	dbAddressIsSet      bool
	CryptoKey           string `env:"CRYPTO_KEY" json:"crypto_key"`
	cryptoKeyIsSet      bool
	ConfigFilePath      string `env:"CONFIG" json:"-"`
	configFilePathIsSet bool
	TokenExpiry         time.Duration `env:"TOKEN_EXPIRY" json:"token_expiry"`
	tokenExpiryIsSet    bool
}

// ServerConfigBuilder is a builder for constructing a ServerConfig instance.
type ServerConfigBuilder struct {
	Config ServerConfig
}

func (c *ServerConfig) SetDefaults() {
	c.Address = defaultAddress
	c.DBAddress = defaultDBAddress
	c.CryptoKey = defaultCryptoKey
	c.ConfigFilePath = defaultConfigFilePath
	c.TokenExpiry = defaultTokenExpiry

}

// WithAddress sets the address in the ServerConfig.
func (c *ServerConfigBuilder) WithAddress(address string) *ServerConfigBuilder {
	c.Config.Address = address
	c.Config.addressIsSet = true
	return c
}

// WithDSN sets the database DSN in the ServerConfig.
func (c *ServerConfigBuilder) WithDSN(dsn string) *ServerConfigBuilder {
	c.Config.DBAddress = dsn
	c.Config.dbAddressIsSet = true
	return c
}

// WithCryptoKey sets the crypto key flag in the ServerConfig.
func (c *ServerConfigBuilder) WithCryptoKey(path string) *ServerConfigBuilder {
	c.Config.CryptoKey = path
	c.Config.cryptoKeyIsSet = true
	return c
}

// WithConfigFile sets the path to JSON configuration file
func (c *ServerConfigBuilder) WithConfigFile(configFilePath string) *ServerConfigBuilder {
	c.Config.ConfigFilePath = configFilePath
	c.Config.configFilePathIsSet = true
	return c
}

// WithTokenExpiry sets the authorization token time to live
func (c *ServerConfigBuilder) WithTokenExpiry(exp time.Duration) *ServerConfigBuilder {
	c.Config.TokenExpiry = exp
	c.Config.tokenExpiryIsSet = true
	return c
}

// FromFile populates the ServerConfig from JSON .
func (c *ServerConfigBuilder) FromFile() *ServerConfigBuilder {

	if c.Config.ConfigFilePath == "" {
		return c
	}
	if err := k.Load(file.Provider(c.Config.ConfigFilePath), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	JSONConfig := ServerConfig{}
	JSONConfig.SetDefaults()

	if err := config.ParseConfigFile(&JSONConfig, c.Config.ConfigFilePath); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	if JSONConfig.Address != defaultAddress && !c.Config.addressIsSet {
		c.WithAddress(JSONConfig.Address)
	}

	if JSONConfig.DBAddress != defaultDBAddress && !c.Config.dbAddressIsSet {
		c.WithDSN(JSONConfig.DBAddress)
	}

	if JSONConfig.CryptoKey != defaultCryptoKey && !c.Config.cryptoKeyIsSet {
		c.WithCryptoKey(JSONConfig.CryptoKey)
	}

	if JSONConfig.TokenExpiry != defaultTokenExpiry && !c.Config.tokenExpiryIsSet {
		c.WithTokenExpiry(JSONConfig.TokenExpiry)
	}

	return c
}

// FromEnv populates the ServerConfig from environment variables.
func (c *ServerConfigBuilder) FromEnv() *ServerConfigBuilder {
	if err := env.Parse(&c.Config); err != nil {
		log.Fatal(err)
	}
	_, addressIsSet := os.LookupEnv("ADDRESS")
	if addressIsSet {
		c.Config.addressIsSet = true
	}
	_, cryptoKeyIsSet := os.LookupEnv("CRYPTO_KEY")
	if cryptoKeyIsSet {
		c.Config.cryptoKeyIsSet = true
	}
	_, configFilePathIsSet := os.LookupEnv("CONFIG")
	if configFilePathIsSet {
		c.Config.configFilePathIsSet = true
	}
	_, dsnSet := os.LookupEnv("DATABASE_DSN")
	if dsnSet {
		c.Config.dbAddressIsSet = true
	}
	return c
}

// GetConfigs returns the fully constructed ServerConfig by combining
// configurations from environment variables and command-line flags.
// It validates the address to ensure it is properly set.
func GetConfigs() (ServerConfig, error) {
	var c ServerConfigBuilder
	c.Config.SetDefaults()
	c.FromEnv().FromFile()
	if !validation.ValidateAddress(c.Config.Address) {
		return ServerConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}
