package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/testsuite"
	"strings"
	"testing"

	"main/internal/libs/golangutilslib"
	"main/internal/testutils"
	"main/internal/testutils/golangutilsmock"

	"github.com/stretchr/testify/assert"
)

type ProcessData struct {
	Name      string
	recursive bool
	expected  int
}

var (
	suite    = testsuite.NewWithAll(beforeEachTests, afterEachTests)
	mockFile = new(golangutilsmock.FileMock)
)

var (
	testDir      string
	processCases []ProcessData
)

func beforeEachTests(t *testing.T) {
	testDir = testutils.GetDirToTest()
	setupCommand()
}

func afterEachTests(t *testing.T) {
	mockFile.AssertExpectations(t)
	mockFile.ExpectedCalls = nil // Clean and reset
}

func processTestProcessWithoutArg(t *testing.T, data ProcessData) {
	currentDir := logic.Ternary(strings.Contains(data.Name, "Empty"), file.ResolvePath(testDir, "dir2"), testDir)
	args := logic.Ternary(data.recursive, []string{"-r"}, []string{})
	golangutilslib.FuncGetCurrentDir = mockFile.GetCurrentDir // Inject Mocks
	mockFile.On("GetCurrentDir").Return(currentDir, nil).Once()
	testutils.ForceCobraReadArgs(t, args...)
	output, err := logic.CaptureOutput(process)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d\n", data.expected), output, fmt.Sprintf("Expected to have: %d", data.expected))
}

func TestProcessWithoutArg(t *testing.T) {
	processCases = []ProcessData{
		{"Non-recursive", false, 3},
		{"Recursive", true, 4},
		{"Non-recursive Empty", false, 0},
		{"Recursive Empty", true, 0},
	}
	// Save original vars value for inject mocks
	origGetCurrentDir := golangutilslib.FuncGetCurrentDir
	defer func() {
		golangutilslib.FuncGetCurrentDir = origGetCurrentDir
	}()
	testsuite.RunTestCases(t, suite, processCases, processTestProcessWithoutArg)
}

func processTestProcessWithArg(t *testing.T, data ProcessData) {
	dir := logic.Ternary(strings.Contains(data.Name, "Empty"), file.ResolvePath(testDir, "dir2"), testDir)
	args := logic.Ternary(data.recursive, []string{"-w", dir, "-r"}, []string{"-w", dir})
	testutils.ForceCobraReadArgs(t, args...)
	output, err := logic.CaptureOutput(process)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d\n", data.expected), output, fmt.Sprintf("Expected to have: %d", data.expected))
}

func TestProcessWithArg(t *testing.T) {
	processCases = []ProcessData{
		{"Non-recursive", false, 3},
		{"Recursive", true, 4},
		{"Non-recursive Empty", false, 0},
		{"Recursive Empty", true, 0},
	}
	// Save original vars value for inject mocks
	origGetCurrentDir := golangutilslib.FuncGetCurrentDir
	defer func() {
		golangutilslib.FuncGetCurrentDir = origGetCurrentDir
	}()
	testsuite.RunTestCases(t, suite, processCases, processTestProcessWithArg)
}

func TestMain(t *testing.T) {
	testutils.TestMain(t, suite, main)
}
