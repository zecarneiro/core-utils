package main

import (
	"fmt"
	"golangutils/pkg/models"
	"golangutils/pkg/testsuite"
	"testing"

	"main/internal/libs/golangutilslib"
	"main/internal/testutils"
	"main/internal/testutils/golangutilsmock"
)

var suite = testsuite.NewWithBeforeEach(beforeEachTests)

func beforeEachTests(t *testing.T) {
	setupCommand()
}

func TestMain(t *testing.T) {
	testutils.TestMain(t, suite, main)
}

func TestProcess(t *testing.T) {
	comandosDebGet := []string{"sudo deb-get update", "sudo deb-get upgrade"}
	mockPlatform := new(golangutilsmock.PlatformMock)
	mockExec := new(golangutilsmock.ExeMock)

	// Inject Mocks
	origIsWindows := golangutilslib.FuncIsWindows
	origExecRealtime := golangutilslib.FuncExecRealTime
	defer func() {
		golangutilslib.FuncIsWindows = origIsWindows
		golangutilslib.FuncExecRealTime = origExecRealtime
	}()

	cases := []struct {
		msg             string
		isWindows       bool
		expectedCommand string
	}{
		{"Windows Case", true, "topgrade --cleanup --allow-root --skip-notify --yes --disable helm uv deb_get"},
		{"Linux Case", false, fmt.Sprintf("%s powershell", "topgrade --cleanup --allow-root --skip-notify --yes --disable helm uv deb_get")},
	}
	for _, tc := range cases {
		t.Run(tc.msg, func(t *testing.T) {
			// Clean and reset
			mockPlatform.ExpectedCalls = nil
			mockExec.ExpectedCalls = nil

			golangutilslib.FuncIsWindows = mockPlatform.IsWindows
			golangutilslib.FuncExecRealTime = mockExec.ExecRealTime

			mockPlatform.On("IsWindows").Return(tc.isWindows).Once()
			mockExec.On("ExecRealTime", models.Command{Cmd: tc.expectedCommand, UseShell: true, Verbose: true}).Return(nil).Once()
			for _, cmdStr := range comandosDebGet {
				expectedCmd := models.Command{
					Cmd:      cmdStr,
					UseShell: true,
					Verbose:  true,
				}
				mockExec.On("ExecRealTime", expectedCmd).Return(nil).Once()
			}
			process()
			mockPlatform.AssertExpectations(t)
			mockExec.AssertExpectations(t)
		})
	}
}
