package checkenv

import "os"

type EnvChecker struct {
	VariableNames []string
}

func NewEnvChecker(variableNames []string) *EnvChecker {
	return &EnvChecker{VariableNames: variableNames}
}

func (e *EnvChecker) CheckAll() bool {
	for _, variableName := range e.VariableNames {
		if _, exists := os.LookupEnv(variableName); !exists {
			return false
		}
	}
	return true
}
