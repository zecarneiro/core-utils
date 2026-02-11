package libs

import (
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"
	"slices"
	"strings"
)

type EnvManager struct {
	EnvFile string
}

func NewEnvManager() *EnvManager {
	return &EnvManager{
		EnvFile: file.JoinPath(system.HomeUserConfigDir(), ".env_common"),
	}
}

func (e *EnvManager) getFileShell(shellType enums.ShellType) string {
	filename := "systemd-env"
	directory := file.JoinPath(system.HomeUserConfigDir(), "shell")
	switch shellType {
	case enums.Bash:
		return file.JoinPath(directory, fmt.Sprintf("%s.sh", filename))
	case enums.PowerShell:
		return file.JoinPath(directory, fmt.Sprintf("%s.ps1", filename))
	case enums.Fish:
		return file.JoinPath(directory, fmt.Sprintf("%s.fish", filename))
	}
	return ""
}

func (e *EnvManager) getFileShellData(shellType enums.ShellType) string {
	fileShell := e.getFileShell(shellType)
	switch shellType {
	case enums.Bash:
		return fmt.Sprintf(`[[ -f "%s" ]] && source "%s"`, fileShell, fileShell)
	case enums.PowerShell:
		return fmt.Sprintf(`if (Test-Path "%s" -PathType Leaf) { . "%s" }`, fileShell, fileShell)
	case enums.Fish:
		return fmt.Sprintf(`test -f "%s"; and source "%s"`, fileShell, fileShell)
	}
	return ""
}

func (e *EnvManager) unsetEnv(name string) {
	if !platform.IsWindows() {
		if file.IsFile(e.EnvFile) {
			match := fmt.Sprintf("^export %s=", name)
			logic.ProcessError(file.DeleteFileLines(e.EnvFile, match, true))
		}
	} else {
		cmdStr := `[Environment]::SetEnvironmentVariable("%s", [NullString]::Value, "[System.EnvironmentVariableTarget]::User")`
		cmd := models.Command{
			Cmd:        fmt.Sprintf(cmdStr, name),
			UseShell:   true,
			ShellToUse: enums.PowerShell,
		}
		logic.ProcessError(exe.ExecRealTime(cmd))
	}
}

func (e *EnvManager) bashTemplate(shellFile string) string {
	return fmt.Sprintf(`
reload-shell() {
	shell_file="%s"
	if [ -f "$shell_file" ]; then
		source "$shell_file"
	fi
}
if [ -f "%s" ]; then
	source "%s"
	eval "$(dircolors)"
fi
`, shellFile, e.EnvFile, e.EnvFile)
}

func (e *EnvManager) fishTemplate(shellFile string) string {
	return fmt.Sprintf(`
function reload-shell
	set shell_file "%s"
	if test -f "$shell_file"
		source "$shell_file"
	end
end
if test -f "%s"
	for line in (cat "%s")
		if test -z "$line"
			continue
		end
		if string match -qr '^\s*#' $line
			continue
		end
		set line (string trim (string replace -r '^export\s+' '' $line))
		set parts (string split -m1 "=" "$line")
		if test (count $parts) -eq 2
			set key $parts[1]
			set value $parts[2]
			switch $key
				case PWD SHLVL _
					continue
			end
			set -gx $key $value
		end
	end
	eval (dircolors -c)
end
`, shellFile, e.EnvFile, e.EnvFile)
}

func (e *EnvManager) powershellTemplate(shellFile string) string {
	return fmt.Sprintf(`
function reload-shell {
	$shell_file = "%s"
	if (Test-Path "$shell_file" -PathType Leaf) {
		. "$shell_file"
	}
}

if (Test-Path "%s") {
	Get-Content "%s" | ForEach-Object {
		$_ = $_.Trim()
		if ([string]::IsNullOrWhiteSpace($_) -or $_ -like '#*') { return }
		$_ = $_ -replace '^export\s+', ''
		$parts = $_ -split '=',2
		if ($parts.Length -eq 2) {
			$key = $parts[0]
			if ($key -notin @('PWD','OLDPWD','SHLVL','_')) {
				$value = $parts[1].Trim('"')
				Set-Item -Path "Env:$key" -Value $value
			}
		}
	}
}
`, shellFile, e.EnvFile, e.EnvFile)
}

func (e *EnvManager) powershellWindowsTemplate(shellFile string) string {
	return fmt.Sprintf(`
function reload-shell {
	$shell_file = "%s"
	if (Test-Path "$shell_file" -PathType Leaf) {
		. "$shell_file"
	}
}
`, shellFile)
}

func (e *EnvManager) bashWindowsTemplate(shellFile string) string {
	return fmt.Sprintf(`
reload-shell() {
	shell_file="%s"
	if [ -f "$shell_file" ]; then
		source "$shell_file"
	fi
}
`, shellFile)
}

func (e *EnvManager) writeToFile(filePath string, data string, isShellProfileFile bool) {
	canWrite := false
	fileWriter := models.FileWriterConfig{
		File:        filePath,
		Data:        data,
		IsAppend:    logic.Ternary(isShellProfileFile, true, false),
		IsCreateDir: true,
		WithUtf8BOM: false,
	}
	if isShellProfileFile {
		if file.IsFile(fileWriter.File) {
			isContainData, _ := file.FileTextContains(fileWriter.File, data, false)
			if !isContainData {
				canWrite = true
			}
		} else {
			canWrite = true
		}
	} else {
		if !file.IsFile(fileWriter.File) {
			canWrite = true
		}
	}
	if canWrite {
		logic.ProcessError(file.WriteFile(fileWriter))
		if !isShellProfileFile {
			logic.ProcessError(exe.Chmod777(fileWriter.File, false))
		}
	}
}

func (e *EnvManager) SetSystemConfig() {
	if !platform.IsWindows() {
		if shell.IsBashInstalled() {
			isAlreadyCreatedShellFile := false
			shellTypeList := []enums.ShellType{enums.Bash, enums.Zsh, enums.Ksh}
			for _, shellType := range shellTypeList {
				shellProfileFile := shell.GetShellProfileFile(shellType)
				templateData := e.bashTemplate(shellProfileFile)
				if !isAlreadyCreatedShellFile {
					isAlreadyCreatedShellFile = true
					e.writeToFile(e.getFileShell(shellType), templateData, false)
				}
				e.writeToFile(shellProfileFile, e.getFileShellData(shellType), true)
			}
		}
		if shell.IsFishInstalled() {
			shellType := enums.Fish
			shellProfileFile := shell.GetShellProfileFile(shellType)
			templateData := e.fishTemplate(shellProfileFile)
			e.writeToFile(e.getFileShell(shellType), templateData, false)
			e.writeToFile(shellProfileFile, e.getFileShellData(shellType), true)
		}
		if shell.IsPowershellInstalled() {
			shellType := enums.PowerShell
			shellProfileFile := shell.GetShellProfileFile(shellType)
			templateData := e.powershellTemplate(shellProfileFile)
			e.writeToFile(e.getFileShell(shellType), templateData, false)
			e.writeToFile(shellProfileFile, e.getFileShellData(shellType), true)
		}
	} else {
		if shell.IsBashInstalled() {
			shellType := enums.Bash
			shellProfileFile := shell.GetShellProfileFile(shellType)
			templateData := e.bashWindowsTemplate(shellProfileFile)
			e.writeToFile(e.getFileShell(shellType), templateData, false)
			e.writeToFile(shellProfileFile, e.getFileShellData(shellType), true)
		}
		if shell.IsPowershellInstalled() {
			shellType := enums.PowerShell
			shellProfileFile := shell.GetShellProfileFile(shellType)
			templateData := e.powershellWindowsTemplate(shellProfileFile)
			e.writeToFile(e.getFileShell(shellType), templateData, false)
			e.writeToFile(shellProfileFile, e.getFileShellData(shellType), true)
		}
	}
}

func (e *EnvManager) RemoveDuplicated(values []string) []string {
	newValues := []string{}
	if len(values) == 0 {
		return values
	}
	for _, val := range values {
		if !slices.Contains(newValues, val) {
			newValues = append(newValues, val)
		}
	}
	return newValues
}

func (e *EnvManager) GetEnvValues(name string) []string {
	goEnv := logic.Ternary(env.Exists(name), env.Get(name), []string{})
	values := []string{}
	if !platform.IsWindows() {
		if file.IsFile(e.EnvFile) {
			err := file.ReadFileLineByLine(e.EnvFile, func(line string) {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					line = strings.TrimPrefix(line, "export ")
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						key := strings.TrimSpace(parts[0])
						valuesStr := strings.TrimSpace(parts[1])
						if key == name {
							valuesStr = strings.Trim(valuesStr, `"`)
							values = append(values, env.ConvValuesArr(valuesStr)...)
						}
					}
				}
			})
			logic.ProcessError(err)
		}
	} else {
		cmdStr := `([Environment]::GetEnvironmentVariable("%s", [System.EnvironmentVariableTarget]::User)) | Out-String`
		cmd := models.Command{
			Cmd:        fmt.Sprintf(cmdStr, name),
			UseShell:   true,
			ShellToUse: enums.PowerShell,
		}
		output, err := exe.Exec(cmd)
		logic.ProcessError(err)
		values = env.ConvValuesArr(output)
	}
	values = e.RemoveDuplicated(append(values, goEnv...))
	return values
}

func (e *EnvManager) Sync(name string) {
	values := e.GetEnvValues(name)
	env.Set(name, values)
}

func (e *EnvManager) UpdateEnv(name string, values []string) {
	if !platform.IsWindows() {
		e.unsetEnv(name)
		if len(values) > 0 {
			valuesStr := env.ConvValuesStr(e.RemoveDuplicated(values))
			fileWriter := models.FileWriterConfig{
				File:        e.EnvFile,
				Data:        fmt.Sprintf(`export %s="%s"`, name, valuesStr),
				IsAppend:    true,
				IsCreateDir: true,
				WithUtf8BOM: false,
			}
			logic.ProcessError(file.WriteFile(fileWriter))
		}
	} else {
		if len(values) == 0 {
			e.unsetEnv(name)
		} else {
			cmdStr := `[Environment]::SetEnvironmentVariable("%s", "%s", [System.EnvironmentVariableTarget]::User)`
			cmd := models.Command{
				Cmd:        fmt.Sprintf(cmdStr, name, env.ConvValuesStr(values)),
				UseShell:   true,
				ShellToUse: enums.PowerShell,
			}
			logic.ProcessError(exe.ExecRealTime(cmd))
		}
	}
}
