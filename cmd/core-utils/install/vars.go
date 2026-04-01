package main

import "main/internal/libs"

var (
	rootDir          string
	systemInstallDir string
	envPathName      string
	envManager       *libs.EnvManager
	installTitleMsg  = "Install Core Utils on the system"
	isUpdateOnly     bool
)
