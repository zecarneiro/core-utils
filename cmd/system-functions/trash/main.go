package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
	"os"
	"time"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var linuxTrashInfoData = `[Trash Info]
Path=%s
DeletionDate=%s`

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "trash <filepath>",
		Short: "Delete file to recycle bin",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func getDevice(path string) string {
	devPath, err := file.GetDevice(path)
	logic.ProcessError(err)
	return devPath
}

func tryGioTrash(path string) bool {
	gioCmd := "gio"
	cmds := []string{gioCmd, "gvfs-trash"}
	for _, cmd := range cmds {
		if _, err := console.Which(cmd); err == nil {
			args := logic.Ternary(cmd == gioCmd, []string{"trash", path}, []string{path})
			_, err := exe.Exec(models.Command{Cmd: cmd, Args: args})
			return err == nil
		}
	}
	return false
}

func processLinux(absPath string) {
	homeDir := system.HomeDir()
	uid := os.Getuid()
	devFile := getDevice(absPath)
	devHome := getDevice(homeDir)
	var trashDir string
	var relativePath string
	if devFile == devHome {
		trashDir = file.JoinPath(system.HomeDir(), ".local", "share", "Trash")
		relPath, err := file.GetRelativePath(absPath, homeDir)
		logic.ProcessError(err)
		relativePath = relPath
	} else {
		mountPoint, err := file.FindMountPoint(absPath)
		logic.ProcessError(err)
		trashDir = file.JoinPath(mountPoint, fmt.Sprintf(".Trash-%d", uid))
		relativePath = absPath
	}
	filesDir := file.JoinPath(trashDir, "files")
	infoDir := file.JoinPath(trashDir, "info")
	logic.ProcessError(file.CreateDirectory(filesDir, true))
	logic.ProcessError(file.CreateDirectory(infoDir, true))
	baseName := file.Basename(absPath)
	destPath := file.JoinPath(filesDir, baseName)
	// Avoid conflicts
	counter := 1
	for {
		if !file.FileExist(destPath) {
			break
		}
		destPath = file.JoinPath(filesDir, fmt.Sprintf("%s_%d", baseName, counter))
		counter++
	}
	// Move
	if file.IsFile(absPath) {
		logic.ProcessError(file.CopyFile(absPath, destPath))
		logic.ProcessError(file.DeleteFile(absPath))
	} else if file.IsDir(absPath) {
		logic.ProcessError(file.CopyDir(absPath, destPath))
		logic.ProcessError(file.DeleteDirectory(absPath))
	} else {
		logic.ProcessError(errors.New("Invalid type of file"))
	}

	// Save trash info
	trashInfoPath := file.JoinPath(infoDir, fmt.Sprintf("%s.trashinfo", file.Basename(destPath)))
	trashInfoData := fmt.Sprintf(linuxTrashInfoData, relativePath, time.Now().Format("2006-01-02T15:04:05"))
	logic.ProcessError(os.WriteFile(trashInfoPath, []byte(trashInfoData), 0o644))
}

func processWindows(absPath string) {
	cmd := models.Command{
		Cmd:      fmt.Sprintf("shell = new-object -comobject \"Shell.Application\"; $shell.Namespace(0).ParseName(\"%s\").InvokeVerb(\"delete\")", absPath),
		UseShell: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func process(filePath string) {
	absPath, err := file.GetFullPath(filePath)
	logic.ProcessError(err)
	if !file.FileExist(absPath) {
		logic.ProcessError(errors.New("given a valid file"))
	}
	if platform.IsWindows() {
		processWindows(absPath)
	} else if platform.IsLinux() {
		if !tryGioTrash(absPath) {
			processLinux(absPath)
		}
	} else {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
}

func main() {
	cobralib.Run()
}
