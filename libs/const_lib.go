package libs

import "golangutils"

var LoggerUtils = golangutils.NewLoggerUtils()
var SystemUtils = golangutils.NewSystemUtils(LoggerUtils)
var ConsoleUtils = golangutils.NewConsoleUtils(LoggerUtils, SystemUtils)
