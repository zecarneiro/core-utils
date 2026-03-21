package golangutilslib

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/platform"
)

var (
	// Platform
	FuncIsWindows = platform.IsWindows

	// Logic

	// Exec
	FuncExecRealTime = exe.ExecRealTime

	// File
	FuncGetCurrentDir = file.GetCurrentDir
	FuncCopyDir       = file.CopyDir
)
