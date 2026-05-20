package enums

import "golangutils/pkg/enums"

const (
	APP_VERSION_KEY         = "__VERSION_APP__"
	DOWNLOADED_FILE_ARG_KEY = "__DOWNLOADED_FILE_ARG__"
)

var ValidShellList = []string{enums.PowerShell.String(), enums.Bash.String()}
