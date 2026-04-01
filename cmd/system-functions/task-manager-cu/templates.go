package main

var (
	windowsTemplate = `@echo off
%s
`
	linuxTemplate = `[Unit]
Description=%s

[Service]
Type=oneshot
%s
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
`
)
