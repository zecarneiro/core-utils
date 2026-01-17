package main

import (
	"fmt"
	"golangutils/pkg/console"
)

func main() {
	fmt.Println(console.GetCurrentShell().String())
}
