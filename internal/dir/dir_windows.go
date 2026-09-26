//go:build windows

package dir

import (
	"golangutils/pkg/file"
	"golangutils/pkg/system"
)

func GetScoopRootDir() string {
	return file.JoinPath(system.HomeDir(), "scoop")
}
