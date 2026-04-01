package main

type MenuEntryData struct {
	name       string
	exec       string
	args       []string
	argsStr    string
	icon       string
	categories string
	comment    string
	terminal   bool
	runAsAdmin bool
}

var menuEntryData = MenuEntryData{}
