package golangutilsmock

import (
	"golangutils/pkg/models"

	"github.com/stretchr/testify/mock"
)

type ExeMock struct {
	mock.Mock
}

func (m *ExeMock) ExecRealTime(cmd models.Command) error {
	args := m.Called(cmd)
	return args.Error(0)
}
