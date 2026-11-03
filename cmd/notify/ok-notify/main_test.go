package main

import (
	"golangutils/pkg/testsuite"
	"testing"

	"main/internal/testutils"
)

var suite = testsuite.NewWithBeforeEach(beforeEachTests)

func beforeEachTests(t *testing.T) {
	setupCommand()
}

func TestMain(t *testing.T) {
	testutils.TestMain(t, suite, main)
}
