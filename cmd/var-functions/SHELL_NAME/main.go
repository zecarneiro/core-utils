package main

import (
	"fmt"
	"main/libs"
	"strings"
)

func main() {
	fmt.Println(strings.ToUpper(libs.ShellUtils.CurrentShell.String()))
}
