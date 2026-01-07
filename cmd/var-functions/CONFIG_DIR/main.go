package main

import (
	"main/libs"
)

func main() {
	dir_lib := libs.NewDirsLib()
	config_dir := dir_lib.GetUserConfig()
	libs.LoggerUtils.Log(config_dir)
}
