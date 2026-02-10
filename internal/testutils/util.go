package testutils

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/testsuite"
	"os"
	"testing"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func GetDirToTest() string {
	path := os.Getenv("ROOT_DIR")
	path = file.ResolvePath(path, "test")
	return path
}

func processTestMain(t *testing.T, data struct {
	Name     string
	err      error
	mainFunc func()
	mainArgs []string
},
) {
	runCalled := false
	cobralib.FuncExecute = func() error {
		if data.err == nil {
			cobralib.FuncRun(&cobra.Command{}, []string{})
		}
		return data.err
	}
	if data.err == nil {
		cobralib.FuncRun = func(cmd *cobra.Command, args []string) {
			runCalled = true
		}
	}
	cobralib.CobraCmd.SetArgs(data.mainArgs)
	output, err := logic.CaptureOutput(data.mainFunc)
	if err != nil {
		t.Fatal(err)
	}

	if data.err == nil {
		if !runCalled {
			t.Errorf("The command would be called by Cobra")
			assert.Empty(t, output, "Expected to have no log")
		}
	} else {
		assert.Contains(t, output, "Detect Exit with code: 1", "Expected to have exit message")
		assert.Contains(t, output, "Got error during the execution", "Expected to have error message")
	}
}

func TestMain(t *testing.T, s *testsuite.Suite, mainFunc func(), mainArgs ...string) {
	// Inject Mocks
	origCobraCmd := cobralib.CobraCmd
	defer func() {
		cobralib.CobraCmd = origCobraCmd
	}()
	cases := []struct {
		Name     string
		err      error
		mainFunc func()
		mainArgs []string
	}{
		{"With Error", fmt.Errorf("Got error during the execution"), mainFunc, mainArgs},
		{"With No Error", nil, mainFunc, mainArgs},
	}
	testsuite.RunTestCases(t, s, cases, processTestMain)
}

func ForceCobraReadArgs(t *testing.T, args ...string) {
	err := cobralib.CobraCmd.ParseFlags(args)
	assert.NoError(t, err)
}
