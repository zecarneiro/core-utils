package main

var (
	linuxTemplate = `[Unit]
Description=%s

[Service]
ExecStart=%s
Restart=always

[Install]
WantedBy=multi-user.target
`
)
