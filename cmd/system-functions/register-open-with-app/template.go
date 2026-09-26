package main

const (
	REG_ADD_TEMPLATE = `Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\Applications\%s]
"FriendlyAppName"="%s"

[HKEY_CLASSES_ROOT\Applications\%s\shell]

[HKEY_CLASSES_ROOT\Applications\%s\shell\open]

[HKEY_CLASSES_ROOT\Applications\%s\shell\open\command]
@="\"%s\" \"%s\""
`
	REG_DELETE_TEMPLATE = `Windows Registry Editor Version 5.00

[-HKEY_CLASSES_ROOT\Applications\%s]
`
)
