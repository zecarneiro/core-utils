package main

import "main/internal/libs"

var (
	appName      = "envc"
	scriptCmd    string
	envManager   *libs.EnvManager
	warnningHelp = `I you add env, will be create a file: %s.
So if change any file manaually, please change/add/delete file according with file name.`
)
