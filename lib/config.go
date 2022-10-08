package lib

import (
	"fmt"
	"os"
)

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

func (c *Config) Populate() error {
	email, ok := os.LookupEnv(emailPassEnvVar)
	if !ok {
		return fmt.Errorf(lookupErr, emailPassEnvVar)
	}
	c.cfg[EmailPass] = email

	host, ok := os.LookupEnv(dbHostEnvVar)
	if !ok {
		return fmt.Errorf(lookupErr, dbHostEnvVar)
	}
	c.cfg[DBHost] = host

	user, ok := os.LookupEnv(dbUserEnvVar)
	if !ok {
		return fmt.Errorf(lookupErr, dbUserEnvVar)
	}
	c.cfg[DBUsername] = user

	dbPass, ok := os.LookupEnv(dbPassEnvVar)
	if !ok {
		return fmt.Errorf(lookupErr, dbPassEnvVar)
	}
	c.cfg[DBPass] = dbPass

	dbPort, ok := os.LookupEnv(dbPortEnvVar)
	if !ok {
		return fmt.Errorf(lookupErr, dbPortEnvVar)
	}
	c.cfg[DBPort] = dbPort

	if _, ok := os.LookupEnv(devEnvVar); !ok {
		c.cfg[Development] = "false"
	} else {
		c.cfg[Development] = "true"
	}

	return nil
}
