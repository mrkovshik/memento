// Package config provides configuration handling for the server, allowing
// configurations to be set via environment variables or command-line flags.
package config

import (
	"errors"
	"log"
	"os"

	"github.com/eschao/config"
	"github.com/eschao/config/env"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mrkovshik/memento/internal/validation"
)

const (
	defaultConfigFilePath = ""
	defaultAddress        = "localhost:8080"
	defaultDBAddress      = ""
	defaultCryptoKey      = "./public_key.pem"
)

var k = koanf.New(".")

// ServerConfig holds the configuration settings for the server.
type ServerConfig struct {
	Address             string `env:"ADDRESS" json:"address"`
	AddressIsSet        bool   `json:"-"`
	DBAddress           string `env:"DATABASE_DSN" json:"database_dsn"`
	DBAddressIsSet      bool   `json:"-"`
	CryptoKey           string `env:"CRYPTO_KEY" json:"crypto_key"`
	CryptoKeyIsSet      bool   `json:"-"`
	ConfigFilePath      string `env:"CONFIG" json:"-"`
	ConfigFilePathIsSet bool   `json:"-"`
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

}

// WithAddress sets the address in the ServerConfig.
func (c *ServerConfigBuilder) WithAddress(address string) *ServerConfigBuilder {
	c.Config.Address = address
	c.Config.AddressIsSet = true
	return c
}

// WithDSN sets the database DSN in the ServerConfig.
func (c *ServerConfigBuilder) WithDSN(dsn string) *ServerConfigBuilder {
	c.Config.DBAddress = dsn
	c.Config.DBAddressIsSet = true
	return c
}

// WithCryptoKey sets the crypto key flag in the ServerConfig.
func (c *ServerConfigBuilder) WithCryptoKey(path string) *ServerConfigBuilder {
	c.Config.CryptoKey = path
	c.Config.CryptoKeyIsSet = true
	return c
}

// WithConfigFile sets the path to JSON configuration file
func (c *ServerConfigBuilder) WithConfigFile(configFilePath string) *ServerConfigBuilder {
	c.Config.ConfigFilePath = configFilePath
	c.Config.ConfigFilePathIsSet = true
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

	if JSONConfig.Address != defaultAddress && !c.Config.AddressIsSet {
		c.WithAddress(JSONConfig.Address)
	}

	if JSONConfig.DBAddress != defaultDBAddress && !c.Config.DBAddressIsSet {
		c.WithDSN(JSONConfig.DBAddress)
	}

	if JSONConfig.CryptoKey != defaultCryptoKey && !c.Config.CryptoKeyIsSet {
		c.WithCryptoKey(JSONConfig.CryptoKey)
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
		c.Config.AddressIsSet = true
	}
	_, cryptoKeyIsSet := os.LookupEnv("CRYPTO_KEY")
	if cryptoKeyIsSet {
		c.Config.CryptoKeyIsSet = true
	}
	_, configFilePathIsSet := os.LookupEnv("CONFIG")
	if configFilePathIsSet {
		c.Config.ConfigFilePathIsSet = true
	}
	_, dsnSet := os.LookupEnv("DATABASE_DSN")
	if dsnSet {
		c.Config.DBAddressIsSet = true
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
