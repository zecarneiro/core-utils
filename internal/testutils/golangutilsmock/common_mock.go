package golangutilsmock

import (
	"golangutils/pkg/testsuite"
)

func EnableRunningTests() {
	testsuite.IsRunningTests = true
}
