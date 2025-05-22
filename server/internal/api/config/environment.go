package config

import "fmt"

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

func NewEnvironment(env string) (Environment, error) {
	if env != string(EnvironmentDevelopment) && env != string(EnvironmentProduction) {
		return "", fmt.Errorf("invalid environment: %s", env)
	}
	return Environment(env), nil
}

func (e Environment) IsDevelopment() bool {
	return e == "development"
}
