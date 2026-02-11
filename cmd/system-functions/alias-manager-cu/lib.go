package main

import (
	"fmt"
	"golangutils/pkg/logic"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"
	"regexp"

	"main/internal/dir"
)

func validateName(name string) {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(name) || !re.MatchString(name) {
		logic.ProcessError(fmt.Errorf("%s: Invalid given name: %s. Name must not be empty and accept only A-Z, a-z, 0-9, - and _", appName, name))
	}
}

func loadAndValidateVars() {
	currentShell = shell.GetCurrentShellSimple()
	// Load others
	aliasDir = dir.CoreUtilsShellScripts()
}
