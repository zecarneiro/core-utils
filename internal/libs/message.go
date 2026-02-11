package libs

import (
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/slice"
)

func GetValidShellMsg(parentName string, validShellList []enums.ShellType) error {
	return fmt.Errorf("%s - Running [%s]-specific command only", parentName, slice.ObjArrayToStringBySep(validShellList, ","))
}
