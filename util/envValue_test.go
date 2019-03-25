package util

import (
	"os"
	"testing"
)

func TestEnvValueWithDefault(t *testing.T) {
	environmentVariableName := "FOO"
	environmentVariableValue := "bar"
	defer func() {
		_ = os.Unsetenv(environmentVariableName)
	}()
	err := os.Setenv(environmentVariableName, environmentVariableValue)
	if err != nil {
		t.Error(err)
		return
	}

	result := EnvValueWithDefault(environmentVariableName, "")
	if result != environmentVariableValue {
		t.Errorf("Expected to get value of environment variable, which is '%s', but got '%s'...", environmentVariableValue, result)
	}
}

func TestEnvValueWithDefault_variableNotExisting(t *testing.T) {
	result := EnvValueWithDefault("NOT_EXISTING", "default")

	expectedValue := "default"
	if result != expectedValue {
		t.Errorf("Expected to get default value, which is '%s', but got '%s'...", expectedValue, result)
	}
}

func TestEnvValue(t *testing.T) {
	environmentVariableName := "FOO"
	environmentVariableValue := "bar"
	defer func() {
		_ = os.Unsetenv(environmentVariableName)
	}()
	err := os.Setenv(environmentVariableName, environmentVariableValue)
	if err != nil {
		t.Error(err)
		return
	}

	result, err := EnvValue(environmentVariableName)
	if err != nil {
		t.Errorf("Didn't expect an error, but got: %v", err)
		return
	}
	if result != environmentVariableValue {
		t.Errorf("Expected to get value of environment variable, which is '%s', but got '%s'...", environmentVariableValue, result)
	}
}

func TestEnvValue_variableNotExisting(t *testing.T) {
	_, err := EnvValue("NOT_EXISTING")

	if err == nil {
		t.Errorf("Expected error as the environment variable isn't existing.")
		return
	}

	_, ok := err.(*EnvValueNotExistingError)
	if !ok {
		t.Errorf("Expected EnvValueNotExistingError, but got: %v", err)
	}
}
