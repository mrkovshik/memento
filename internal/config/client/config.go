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
	defaultServerCert     = `-----BEGIN CERTIFICATE-----
MIIDDzCCAfegAwIBAgIUKavEmDjwTwkRvtSr97Awv78P448wDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI0MDkwMTAzMjA0MVoXDTI1MDkw
MTAzMjA0MVowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAoAB0eRizn4qNbIeJ0usx4TIwkg1lkckUeNEJVGHsrJ5E
JN4+kT/ueV1wFDZffE/2kbiaA3fbO4y6dBSvFG3nLwLRozx4/ovtDaLJfLwwMhAz
xkyB1OW1QX460HuIUeWkmXpbOaxzfRqlzq/F3SxHv4GS72yXkbrTjDOq4cLEJ+CJ
HoeD6XU/HhOE6mrv5yWyU4qwtbTm+lwhx8yMF7Q1J89QqRLh2wujdck23txJrswz
kzclgqfWXU3mfpJ0E5tvOSNj055rusw85VJvB6Z+1rLjHI0w+wEfiWi9LQoHYjGD
h2k0EovXTSVyvMFZ5vgkNqiYgNgo9DjnuN5oCLI99wIDAQABo1kwVzALBgNVHQ8E
BAMCBDAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwFAYDVR0RBA0wC4IJbG9jYWxob3N0
MB0GA1UdDgQWBBQTi1N3hUDOkRwGlZyjFkumqP13zzANBgkqhkiG9w0BAQsFAAOC
AQEAWb89u2auAUQf4Rbx8in5z8cOq8umHoKAcn+CE2VgFLSjwHKi2NzS5RSpCsQt
iwLC1Mmhtc6KVeOjuFa4XvS3Ken75tG8byrP74WvzupAJHHix/XyukU9kUuwjnmL
ij7W1iavWrUPW+VX/PbNr/eIKCct6O5YGMN7dzkRv60YYXX81CZXxj0n5qV5IG0p
3IX4+Rdv2JqX4ZCCBA4N2hoEYpgsK+lyB8HBSG6jC/Dgl104mWkI98IEX0W4dZbz
dHU/bbAPoqYQz74ub0Cxks9Qo6k4cEtCqqOH28xVk7AQmiWICHyhBt/hemSEzxxP
Fujpr7VIiuf8sj3X4rQPqYqi3g==
-----END CERTIFICATE-----`
)

var k = koanf.New(".")

// ClientConfig holds the configuration settings for the client.
type ClientConfig struct {
	Address                string `env:"ADDRESS" json:"address"`
	AddressIsSet           bool   `json:"-"`
	CryptoKey              string `env:"CRYPTO_KEY" json:"crypto_key"`
	CryptoKeyIsSet         bool   `json:"-"`
	ConfigFilePath         string `env:"CONFIG" json:"-"`
	ConfigFilePathIsSet    bool   `json:"-"`
	ServerCertificate      string `env:"SERVER_CERT" json:"server_certificate"`
	ServerCertificateIsSet bool   `json:"-"`
}

// ClientConfigBuilder is a builder for constructing a ClientConfig instance.
type ClientConfigBuilder struct {
	Config ClientConfig
}

func (c *ClientConfig) SetDefaults() {
	c.Address = defaultAddress
	c.CryptoKey = defaultCryptoKey
	c.ConfigFilePath = defaultConfigFilePath
	c.ServerCertificate = defaultServerCert

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

// WithServerCertificate sets TLS certificate
func (c *ClientConfigBuilder) WithServerCertificate(serverCertificate string) *ClientConfigBuilder {
	c.Config.ServerCertificate = serverCertificate
	c.Config.ServerCertificateIsSet = true
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

	if JSONConfig.ServerCertificate != defaultServerCert && !c.Config.ServerCertificateIsSet {
		c.WithServerCertificate(JSONConfig.ServerCertificate)
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
	_, serverCertIsSet := os.LookupEnv("SERVER_CERT")
	if serverCertIsSet {
		c.Config.ServerCertificateIsSet = true
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
