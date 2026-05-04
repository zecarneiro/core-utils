package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"slices"
	"strings"

	"main/cmd/tools-functions/dev-container-generator/models"
	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"

	modelsext "golangutils/pkg/models"
)

var (
	typeArg             string
	onCreateCmdArg      string
	updateContentCmdArg string
	postCreateCmdArg    string
	postStartCmdArg     string
	portsArg            string
	vscodeArg           string
	jetbrainsArg        string
	languageVersionArg  string

	devContainerDir string

	validTypeArg           = []string{"go"}
	defaultLanguageVersion = map[string]string{
		"go": "1.24.11",
	}
	workDirArg       = "/workspace"
	ubuntuVersionArg = "24.04"
	remoteUserArg    = "ubuntu"
	containerUserArg = "ubuntu"
)

func init() {
	console.EnableFeatures()
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "dev-container-generator",
		Short: "Set full permission for a file or directory",
	}
	cobralib.CobraCmd.Flags().StringVarP(&typeArg, "type", "t", "", fmt.Sprintf("Type of dev container to generate. Accept only: %s", slice.ArrayToString(validTypeArg)))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("type"))

	// Lifecycle commands
	cobralib.CobraCmd.Flags().StringVarP(&onCreateCmdArg, "on-create-cmd", "c", "", "Command to run on create container")
	cobralib.CobraCmd.Flags().StringVarP(&updateContentCmdArg, "update-content-cmd", "u", "", "Command to run to update content in container")
	cobralib.CobraCmd.Flags().StringVarP(&postCreateCmdArg, "post-create-cmd", "p", "", "Command to run post create container")
	cobralib.CobraCmd.Flags().StringVarP(&postStartCmdArg, "post-start-cmd", "s", "", "Command to run post start container")

	// Users
	cobralib.CobraCmd.Flags().StringVarP(&remoteUserArg, "remote-user", "r", remoteUserArg, "Remote User")
	cobralib.CobraCmd.Flags().StringVarP(&containerUserArg, "container-user", "U", containerUserArg, "Container User")

	// Files/config
	cobralib.CobraCmd.Flags().StringVarP(&portsArg, "ports", "P", "", "JSON file with ports config, i.e. forwardPorts and portsAttributes")
	cobralib.CobraCmd.Flags().StringVarP(&workDirArg, "workdir", "w", workDirArg, "Work directory")
	cobralib.CobraCmd.Flags().StringVarP(&vscodeArg, "vscode", "v", "", "JSON file with vscode config, i.e. settings and extensions")
	cobralib.CobraCmd.Flags().StringVarP(&jetbrainsArg, "jetbrains", "j", "", "JSON file with jetbrains config, i.e. settings and plugins")

	// Versions
	cobralib.CobraCmd.Flags().StringVarP(&languageVersionArg, "language-version", "l", "", fmt.Sprintf("Version of language. Default by type: %s", slice.MapToString(defaultLanguageVersion)))
	cobralib.CobraCmd.Flags().StringVarP(&ubuntuVersionArg, "ubuntu-version", "b", ubuntuVersionArg, "Version of Ubuntu SO")

	cobralib.WithRun(process)
}

func buildPorts(containerData *models.MDevContainer) {
	portsArg = file.ResolvePath(portsArg)
	if file.IsFile(portsArg) {
		jsonData, err := file.ReadJsonFile[models.MDevContainer](portsArg)
		logic.ProcessError(err)
		containerData.ForwardPorts = jsonData.ForwardPorts
		containerData.PortsAttributes = jsonData.PortsAttributes
	}
}

func buildVscode() models.MVScode {
	terminalIntegratedKey := "terminal.integrated.cwd"
	vscodeArg = file.ResolvePath(vscodeArg)
	if file.IsFile(vscodeArg) {
		jsonData, err := file.ReadJsonFile[models.MVScode](vscodeArg)
		logic.ProcessError(err)
		if !slice.MapExistKey(jsonData.Settings, terminalIntegratedKey) {
			jsonData.Settings[terminalIntegratedKey] = workDirArg
		}
		return jsonData
	}
	return models.MVScode{Settings: map[string]any{terminalIntegratedKey: workDirArg}}
}

func buildJetBrains() models.MJetBrains {
	jetbrainsArg = file.ResolvePath(jetbrainsArg)
	if file.IsFile(jetbrainsArg) {
		jsonData, err := file.ReadJsonFile[models.MJetBrains](jetbrainsArg)
		logic.ProcessError(err)
		return jsonData
	}
	return models.MJetBrains{}
}

func fillDevContainerDir() {
	dirBasename := ".devcontainer"
	currentDir := libs.GetCurrentDir(true)
	dirname := file.Basename(currentDir)
	if dirname == dirBasename {
		devContainerDir = currentDir
	} else {
		devContainerDir = file.JoinPath(currentDir, dirBasename)
	}
}

func buildJsonFile() *models.MDevContainer {
	devContainerData := models.MDevContainer{
		Name: fmt.Sprintf(`%s Dev Container`, strings.ToUpper(typeArg)),
		Build: models.MBuild{
			Dockerfile: "Dockerfile",
			Context:    "..",
			Args: map[string]string{
				"SO_VARIANT_ARG": ubuntuVersionArg,
				"WORK_DIR_ARG":   workDirArg,
			},
		},
		OnCreateCommand:      onCreateCmdArg,
		UpdateContentCommand: updateContentCmdArg,
		PostCreateCommand:    postCreateCmdArg,
		PostStartCommand:     postStartCmdArg,
		RemoteUser:           remoteUserArg,
		ContainerUser:        containerUserArg,
		Privileged:           false,
		WorkspaceFolder:      workDirArg,
		WorkspaceMount:       fmt.Sprintf(`source=${localWorkspaceFolder},target=%s,type=bind,consistency=cached`, workDirArg),
		Customizations: models.MCustomizations{
			VScode:    buildVscode(),
			JetBrains: buildJetBrains(),
		},
	}
	buildPorts(&devContainerData)
	return &devContainerData
}

func buildDockerFile(devContainerData *models.MDevContainer) string {
	if typeArg == "go" {
		devContainerData.Build.Args["GO_VERSION_ARG"] = logic.Ternary(str.IsEmpty(languageVersionArg), defaultLanguageVersion[typeArg], languageVersionArg)
		return fmt.Sprintf(dockerfileTemplate, dockerfileGoTemplate)
	}
	return ""
}

func writeFiles(devContainerData *models.MDevContainer, dockerData string) {
	jsonBasename := "devcontainer"
	dockerBasename := "Dockerfile"
	devContainerJson := file.JoinPath(devContainerDir, fmt.Sprintf("%s-generated.json", jsonBasename))
	dockerFile := file.JoinPath(devContainerDir, fmt.Sprintf("%s-generated", dockerBasename))
	logic.ProcessError(file.WriteJsonFile(devContainerJson, devContainerData, true))
	logic.ProcessError(file.WriteFile(modelsext.FileWriterConfig{File: dockerFile, Data: dockerData, IsAppend: false, IsCreateDir: true, WithUtf8BOM: false}))
	logger.Ok("Created files:")
	logger.Prompt(devContainerJson)
	logger.Prompt(dockerFile)
}

func process() {
	if !slices.Contains(validTypeArg, typeArg) {
		logic.ProcessError(fmt.Errorf("Invalid type. Only accept %s", slice.ArrayToString(validTypeArg)))
	}
	fillDevContainerDir()
	devContainerData := buildJsonFile()
	dockerData := buildDockerFile(devContainerData)
	writeFiles(devContainerData, dockerData)
}

func main() {
	cobralib.Run()
}
