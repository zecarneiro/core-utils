package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"strings"
)

func getNameNormalized(name string, toCamelCase bool) string {
	if toCamelCase {
		return str.ToCamelCase(name, true)
	}
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func printOkMsg(entryFilepath string) {
	logger.Ok(fmt.Sprintf("Created Menu Entry: %s", entryFilepath))
}

func writeToFile(filePath string, data string) {
	fileConfigAdmin := models.FileWriterConfig{
		File:        filePath,
		Data:        data,
		IsAppend:    false,
		WithUtf8BOM: false,
	}
	logic.ProcessError(file.WriteFile(fileConfigAdmin))
}
