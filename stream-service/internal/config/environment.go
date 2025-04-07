package config

import "fmt"

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

func NewEnvironment(env string) (Environment, error) {
	if env != string(EnvDevelopment) && env != string(EnvProduction) {
		return "", fmt.Errorf("invalid environment: %s", env)
	}
	return Environment(env), nil
}

func (e Environment) IsDevelopment() bool {
	return e == "development"
}
