package args

import (
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"slices"

	cmdenums "main/cmd/package-functions/installer-by-url/enums"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

type BaseArgs struct {
	AppName         string
	Url             string
	FileExtension   string
	ScriptOrCommand string
	Shell           string
	IsCommand       bool
	RunAfterInit    bool
}

var (
	appNameFlags = cobralib.FlagsOptions[string]{
		Name:       "name",
		Shorthand:  "n",
		Usage:      "Name of the APP",
		IsRequired: true,
		Default:    "",
	}
	urlFlags = cobralib.FlagsOptions[string]{
		Name:       "url",
		Shorthand:  "u",
		Usage:      fmt.Sprintf("URL to download file. I URL contain version to upgrade app, please, pass %s and will be replaced by new getting version.", cmdenums.APP_VERSION_KEY),
		IsRequired: true,
		Default:    "",
	}
	fileExtensionFlags = cobralib.FlagsOptions[string]{
		Name:       "file-extension",
		Shorthand:  "e",
		Usage:      "File extension(Without dot(.)), like: deb, exe, etc.",
		IsRequired: false,
		Default:    "",
	}
	scriptOrCommandFlags = cobralib.FlagsOptions[string]{
		Name:       "script-or-command",
		Shorthand:  "C",
		Usage:      fmt.Sprintf("Shell script to call after to download the file from url. This script must receive the 1st argument the downloaded file. If you use as command, you may need this KEY(%s) to pass file to command.", cmdenums.DOWNLOADED_FILE_ARG_KEY),
		IsRequired: true,
		Default:    "",
	}
	shellFlags = cobralib.FlagsOptions[string]{
		Name:      "shell",
		Shorthand: "s",
		Usage:     fmt.Sprintf("Shell to run the script. Only accept: %s", slice.ArrayToStringBySep(cmdenums.ValidShellList, ",")),
		Default:   logic.Ternary(platform.IsWindows(), enums.PowerShell.String(), enums.Bash.String()),
	}
	isCommandFlags = cobralib.FlagsOptions[bool]{
		Name:      "is-command",
		Shorthand: "c",
		Usage:     "If set, will inform that, script-or-command arg is command",
		Default:   false,
	}
	runAfterInitFlags = cobralib.FlagsOptions[bool]{
		Name:      "run-after-init",
		Shorthand: "R",
		Usage:     "If set, will run after init config",
		Default:   false,
	}
)

func NewBaseArgs() *BaseArgs {
	return &BaseArgs{}
}

func buildBaseFlags(inputCmd *cobra.Command) {
	// App name
	inputCmd.Flags().StringP(appNameFlags.Name, appNameFlags.Shorthand, appNameFlags.Default, appNameFlags.Usage)
	logic.ProcessError(inputCmd.MarkFlagRequired(appNameFlags.Name))
	// URL
	inputCmd.Flags().StringP(urlFlags.Name, urlFlags.Shorthand, urlFlags.Default, urlFlags.Usage)
	logic.ProcessError(inputCmd.MarkFlagRequired(urlFlags.Name))
	// Downloaded file extension
	inputCmd.Flags().StringP(fileExtensionFlags.Name, fileExtensionFlags.Shorthand, fileExtensionFlags.Default, fileExtensionFlags.Usage)
	// Shell Script or Command
	inputCmd.Flags().StringP(scriptOrCommandFlags.Name, scriptOrCommandFlags.Shorthand, scriptOrCommandFlags.Default, scriptOrCommandFlags.Usage)
	logic.ProcessError(inputCmd.MarkFlagRequired(scriptOrCommandFlags.Name))
	// Shell and Is Command flags
	inputCmd.Flags().StringP(shellFlags.Name, shellFlags.Shorthand, shellFlags.Default, shellFlags.Usage)
	inputCmd.Flags().BoolP(isCommandFlags.Name, isCommandFlags.Shorthand, isCommandFlags.Default, isCommandFlags.Usage)
	inputCmd.Flags().BoolP(runAfterInitFlags.Name, runAfterInitFlags.Shorthand, runAfterInitFlags.Default, runAfterInitFlags.Usage)
}

func (b *BaseArgs) validateBase() {
	shellType := enums.GetShellTypeFromValue(b.Shell)
	if str.IsEmpty(b.AppName) {
		logic.ProcessError(fmt.Errorf("Invalid given APP name"))
	}
	if str.Contains(b.AppName, " ", false) {
		logic.ProcessError(fmt.Errorf("APP name must not contains empty space"))
	}
	if str.IsEmpty(b.FileExtension) {
		logic.ProcessError(fmt.Errorf("Invalid given downloaded file extension"))
	}
	if str.IsEmpty(b.Url) {
		logic.ProcessError(fmt.Errorf("Invalid given URL"))
	}
	if !slices.Contains(cmdenums.ValidShellList, b.Shell) || shellType == enums.UnknownShell {
		logic.ProcessError(fmt.Errorf("Invalid given shell"))
	}
	if !b.IsCommand && !file.IsFile(b.ScriptOrCommand) {
		logic.ProcessError(fmt.Errorf("Invalid given shell script"))
	}
}

func (b *BaseArgs) fillBaseValues(cmd *cobra.Command) {
	var err error
	// AppName
	b.AppName, err = cmd.Flags().GetString(appNameFlags.Name)
	logic.ProcessError(err)
	// URL
	b.Url, err = cmd.Flags().GetString(urlFlags.Name)
	logic.ProcessError(err)
	// file extension
	b.FileExtension, err = cmd.Flags().GetString(fileExtensionFlags.Name)
	logic.ProcessError(err)
	// Shell Script or Command
	b.ScriptOrCommand, err = cmd.Flags().GetString(scriptOrCommandFlags.Name)
	logic.ProcessError(err)
	// Shell
	b.Shell, err = cmd.Flags().GetString(shellFlags.Name)
	logic.ProcessError(err)
	// Is Command
	b.IsCommand, err = cmd.Flags().GetBool(isCommandFlags.Name)
	logic.ProcessError(err)
	// Run after init config
	b.RunAfterInit, err = cmd.Flags().GetBool(runAfterInitFlags.Name)
	logic.ProcessError(err)
}
