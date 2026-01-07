package main

import "main/libs"

func getBBB() string {
	return "BBB"
}

func main() {
	dir_lib := libs.NewDirsLib()
	config_dir := dir_lib.GetUserOpt()
	libs.LoggerUtils.Log(config_dir)
}
