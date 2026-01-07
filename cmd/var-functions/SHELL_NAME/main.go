package main

import (
	"fmt"
	"golangutils/pkg/shell"
)

func main() {
	fmt.Println(shell.GetCurrentShell().String())
}
