package main

var (
	windowsTemplate = `@echo off
echo %s
%s
`
	linuxTemplate = `[Unit]
Description=%s

[Service]
ExecStart=%s
Restart=always
StandardOutput=append:/var/log/%s.log
StandardError=append:/var/log/%s.log

[Install]
WantedBy=multi-user.target
`
)
