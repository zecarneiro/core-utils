package main

import (
	"main/libs"
)

func main() {
	libs.LoggerUtils.Log(libs.SystemUtils.OSName())
}
