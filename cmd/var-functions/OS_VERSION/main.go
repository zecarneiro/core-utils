package main

import (
	"fmt"
	"golangutils/pkg/system"
	"strings"
)

func main() {
	fmt.Println(strings.ToUpper(system.OSVersion()))
}
