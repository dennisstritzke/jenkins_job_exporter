package util

import (
	"fmt"
	"os"
)

func EnvValueWithDefault(environmentVariableName string, defaultValue string) string {
	variableContent := os.Getenv(environmentVariableName)
	if len(variableContent) == 0 {
		return defaultValue
	}
	return variableContent
}

func EnvValue(environmentVariableName string) (string, error) {
	variableContent := os.Getenv(environmentVariableName)
	if len(variableContent) == 0 {
		return "", &EnvValueNotExistingError{
			EnvVarName: environmentVariableName,
		}
	}
	return variableContent, nil
}

type EnvValueNotExistingError struct {
	EnvVarName string
}

func (e *EnvValueNotExistingError) Error() string {
	return fmt.Sprintf("environment variable '%s' empty or not existing.", e.EnvVarName)
}
