package main

import (
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"slices"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	policyArg string
	scopeArg  string

	validPolicies = []string{
		"Unrestricted",
		"RemoteSigned",
		"AllSigned",
		"Restricted",
		"Bypass",
	}
	validScopes = []string{
		"CurrentUser",
		"LocalMachine",
		"Process",
	}
)

func init() { setupCommand() }

func setupCommand() {
	policyDefault := validPolicies[4]
	scopeDefault := validScopes[0]
	cobralib.CobraCmd = &cobra.Command{
		Use:   "set-ps-policy",
		Short: "Set PowerShell execution policy",
		Long:  "Sets the PowerShell execution policy for current user or machine",
	}
	cobralib.CobraCmd.Flags().StringVarP(&policyArg, "policy", "p", policyDefault, fmt.Sprintf("Policy: %s (Default: %s)", slice.ArrayToStringBySep(validPolicies, ", "), policyDefault))
	cobralib.CobraCmd.Flags().StringVarP(&scopeArg, "scope", "s", scopeDefault, fmt.Sprintf("Scope: %s (Default: %s)", slice.ArrayToStringBySep(validScopes, ", "), scopeDefault))
	cobralib.WithRun(process)
}

func process() {
	if !slices.Contains(validPolicies, policyArg) {
		logic.ProcessError(fmt.Errorf("Invalid policy. Valid: %s", slice.ArrayToStringBySep(validPolicies, ", ")))
	}
	if !slices.Contains(validScopes, scopeArg) {
		logic.ProcessError(fmt.Errorf("Invalid scope. Valid: %s", slice.ArrayToStringBySep(validScopes, ", ")))
	}
	//
	cmd := fmt.Sprintf("Set-ExecutionPolicy -ExecutionPolicy %s -Scope %s -Force", policyArg, scopeArg)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, UseShell: true, ShellToUse: enums.PowerShell}))
	logger.Ok("Set policy to " + policyArg + " for " + scopeArg)
}

func main() {
	cobralib.Run()
}
