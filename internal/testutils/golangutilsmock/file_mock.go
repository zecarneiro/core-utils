package golangutilsmock

import (
	"github.com/stretchr/testify/mock"
)

type FileMock struct {
	mock.Mock
}

func (m *FileMock) GetCurrentDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
