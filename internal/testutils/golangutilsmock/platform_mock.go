package golangutilsmock

import (
	"github.com/stretchr/testify/mock"
)

type PlatformMock struct {
	mock.Mock
}

func (m *PlatformMock) IsWindows() bool {
	args := m.Called()
	return args.Bool(0)
}
