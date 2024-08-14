// Package config provides configuration handling for the client, allowing
// configurations to be set via environment variables.
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
	defaultAddress        = "localhost:3200"
	defaultCryptoKey      = "./private_key.pem"
)

var k = koanf.New(".")

// ClientConfig holds the configuration settings for the client.
type ClientConfig struct {
	Address             string `env:"ADDRESS" json:"address"`
	AddressIsSet        bool   `json:"-"`
	CryptoKey           string `env:"CRYPTO_KEY" json:"crypto_key"`
	CryptoKeyIsSet      bool   `json:"-"`
	ConfigFilePath      string `env:"CONFIG" json:"-"`
	ConfigFilePathIsSet bool   `json:"-"`
}

// ClientConfigBuilder is a builder for constructing a ClientConfig instance.
type ClientConfigBuilder struct {
	Config ClientConfig
}

func (c *ClientConfig) SetDefaults() {
	c.Address = defaultAddress
	c.CryptoKey = defaultCryptoKey
	c.ConfigFilePath = defaultConfigFilePath

}

// WithAddress sets the address in the ClientConfig.
func (c *ClientConfigBuilder) WithAddress(address string) *ClientConfigBuilder {
	c.Config.Address = address
	c.Config.AddressIsSet = true
	return c
}

// WithCryptoKey sets the crypto key flag in the ClientConfig.
func (c *ClientConfigBuilder) WithCryptoKey(path string) *ClientConfigBuilder {
	c.Config.CryptoKey = path
	c.Config.CryptoKeyIsSet = true
	return c
}

// WithConfigFile sets the path to JSON configuration file
func (c *ClientConfigBuilder) WithConfigFile(configFilePath string) *ClientConfigBuilder {
	c.Config.ConfigFilePath = configFilePath
	c.Config.ConfigFilePathIsSet = true
	return c
}

// FromFile populates the ClientConfig from JSON .
func (c *ClientConfigBuilder) FromFile() *ClientConfigBuilder {

	if c.Config.ConfigFilePath == "" {
		return c
	}
	if err := k.Load(file.Provider(c.Config.ConfigFilePath), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	JSONConfig := ClientConfig{}
	JSONConfig.SetDefaults()

	if err := config.ParseConfigFile(&JSONConfig, c.Config.ConfigFilePath); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	if JSONConfig.Address != defaultAddress && !c.Config.AddressIsSet {
		c.WithAddress(JSONConfig.Address)
	}

	if JSONConfig.CryptoKey != defaultCryptoKey && !c.Config.CryptoKeyIsSet {
		c.WithCryptoKey(JSONConfig.CryptoKey)
	}

	return c
}

// FromEnv populates the ClientConfig from environment variables.
func (c *ClientConfigBuilder) FromEnv() *ClientConfigBuilder {
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
	return c
}

// GetConfigs returns the fully constructed ClientConfig from environment variables.
// It validates the address to ensure it is properly set.
func GetConfigs() (ClientConfig, error) {
	var c ClientConfigBuilder
	c.Config.SetDefaults()
	c.FromEnv().FromFile()
	if !validation.ValidateAddress(c.Config.Address) {
		return ClientConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}
