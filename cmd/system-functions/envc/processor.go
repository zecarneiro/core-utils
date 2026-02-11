package main

import (
	"fmt"
	"golangutils/pkg/conv"
	"golangutils/pkg/env"
	"golangutils/pkg/logic"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"regexp"
	"slices"
)

func validateName(name string) {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(name) || !re.MatchString(name) {
		logic.ProcessError(fmt.Errorf("%s: Invalid given name: %s. Name must not be empty and accept only A-Z, a-z, 0-9, - and _", appName, name))
	}
}

func validateValue(values []string, isOnlyOne bool) {
	values = slice.FilterArray(values, func(val string) bool {
		return !str.IsEmpty(val)
	})
	if len(values) == 0 {
		logic.ProcessError(fmt.Errorf("Invalid given value!"))
	}
	if isOnlyOne && len(values) != 1 {
		logic.ProcessError(fmt.Errorf("Not accept multiple values to validade(only 1 value)"))
	}
}

func existsValue(name string, value string) {
	valueArr := env.ConvValuesArr(value)
	validateName(name)
	validateValue(valueArr, true)
	values := envManager.GetEnvValues(name)
	status := slices.Contains(values, value)
	fmt.Println(conv.BoolToString(status))
}

func exists(name string) {
	validateName(name)
	envManager.Sync(name)
	status := env.Exists(name)
	fmt.Println(conv.BoolToString(status))
}

func processAddArr(name string, valuesArg []string, isReplaceAll bool) {
	var values []string
	validateName(name)
	validateValue(valuesArg, false)
	if isReplaceAll {
		values = valuesArg
	} else {
		values = append(valuesArg, envManager.GetEnvValues(name)...)
	}
	envManager.UpdateEnv(name, values)
}

func processClean(name string) {
	validateName(name)
	envManager.Sync(name)
	if env.Exists(name) {
		values := envManager.GetEnvValues(name)
		envManager.UpdateEnv(name, envManager.RemoveDuplicated(values))
	}
}

func processDelete(name string, valuesArg string) {
	validateName(name)
	envManager.Sync(name)
	if env.Exists(name) {
		valuesArgArr := env.ConvValuesArr(valuesArg)
		if len(valuesArgArr) == 0 {
			envManager.UpdateEnv(name, valuesArgArr)
		} else {
			validateValue(valuesArgArr, true)
			valuesArgArr := env.ConvValuesArr(valuesArg)
			values := envManager.GetEnvValues(name)
			filteredValues := slice.FilterArray(values, func(val string) bool {
				return !slices.Contains(valuesArgArr, val)
			})
			if len(filteredValues) == 0 {
				processDelete(name, "")
			} else if len(values) != len(filteredValues) {
				processAddArr(name, filteredValues, true)
			}
		}
	}
}
