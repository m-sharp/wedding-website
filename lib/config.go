package lib

import (
	"fmt"
	"os"
)

// ToDo: This is probably complex enough to warrant a file now...
const (
	emailPassEnvVar = "EMAILPASSWORD"
	EmailPass       = "Email"

	dbHostEnvVar = "DBHOST"
	DBHost       = "Host"

	dbUserEnvVar = "DBUSER"
	DBUsername   = "User"

	dbPassEnvVar = "DBPASSWORD"
	DBPass       = "Pass"

	dbPortEnvVar = "DBPORT"
	DBPort       = "Port"

	devEnvVar   = "DEV"
	Development = "Development"

	webAdminEnvVar = "WEBUSER"
	WebAdminUser   = "WebAdminUser"

	webAdminPassEnvVar = "WEBPASS"
	WebAdminPass       = "WebAdminPass"

	recaptchaSecretEnvVar = "RECAPTCHASEC"
	RecaptchaSecret       = "RecaptchaSecret"

	csrfSecretEnvVar = "CSRFSEC"
	CSRFSecret       = "CSRFSecret"

	lookupErr = "ENVVAR for %q not found"
)

type Config struct {
	cfg map[string]string
}

func NewConfig() (*Config, error) {
	c := &Config{
		cfg: map[string]string{},
	}

	if err := c.Populate(); err != nil {
		return c, err
	}

	return c, nil
}

func (c *Config) Get(key string) (string, error) {
	value, ok := c.cfg[key]
	if !ok {
		return "", fmt.Errorf("config key %q not found", key)
	}
	return value, nil
}

func (c *Config) Set(key, value string) {
	c.cfg[key] = value
}

var (
	lookupMap = map[string]string{
		emailPassEnvVar:       EmailPass,
		dbHostEnvVar:          DBHost,
		dbUserEnvVar:          DBUsername,
		dbPassEnvVar:          DBPass,
		dbPortEnvVar:          DBPort,
		webAdminEnvVar:        WebAdminUser,
		webAdminPassEnvVar:    WebAdminPass,
		recaptchaSecretEnvVar: RecaptchaSecret,
		csrfSecretEnvVar:      CSRFSecret,
	}
)

func (c *Config) Populate() error {
	for envVarKey, cfgKey := range lookupMap {
		val, ok := os.LookupEnv(envVarKey)
		if !ok {
			return fmt.Errorf(lookupErr, envVarKey)
		}
		c.cfg[cfgKey] = val
	}

	if _, ok := os.LookupEnv(devEnvVar); !ok {
		c.cfg[Development] = "false"
	} else {
		c.cfg[Development] = "true"
	}

	return nil
}
