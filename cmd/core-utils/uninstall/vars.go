package main

import "main/internal/libs"

var (
	envManager        *libs.EnvManager
	systemInstallDir  string
	envPathName       string
	uninstallTitleMsg = "Uninstall Core Utils on the system"
)
