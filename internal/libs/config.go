package libs

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/slice"
	"slices"

	"main/internal/dir"
)

type Config struct {
	PromptStyle          int      `json:"promptStyle,omitempty"`
	ScriptPackageManager []string `json:"scriptPackageManager,omitempty"`
	configFile           string
	availablePromptStyle []int
}

func NewConfig() *Config {
	config := &Config{}
	config.Read()
	return config
}

func (c *Config) loadConfig() {
	if len(c.configFile) == 0 {
		c.configFile = file.JoinPath(dir.CoreUtilsUserConfig(), "config.json")
	}
	c.availablePromptStyle = []int{1, 2, 3, 4}
}

func (c *Config) loadDefault() {
	if !c.IsValidPromptStyle(c.PromptStyle) {
		c.PromptStyle = 4
	}
	if c.ScriptPackageManager == nil {
		c.ScriptPackageManager = []string{}
	}
	c.loadConfig()
}

func (c *Config) IsValidPromptStyle(value int) bool {
	return slices.Contains(c.availablePromptStyle, value)
}

func (c *Config) Write() {
	if !c.IsValidPromptStyle(c.PromptStyle) {
		logic.ProcessError(fmt.Errorf("available PromptStyle: %s", slice.ObjArrayToString(c.availablePromptStyle)))
	}
	c.loadConfig()
	logic.ProcessError(file.WriteJsonFile(c.configFile, c, false))
}

func (c *Config) Read() {
	c.loadConfig()
	configData, err := file.ReadJsonFile[Config](c.configFile)
	if err != nil {
		c.loadDefault()
		if file.IsFile(c.configFile) {
			logic.ProcessError(fmt.Errorf("Erro ao ler config (usando padrão): %v", err))
		}
		c.Write()
	} else {
		*c = configData
		c.loadConfig()
	}
}
