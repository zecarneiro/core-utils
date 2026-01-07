package libs

import "golangutils"

type DirsLib struct{}

func NewDirsLib() *DirsLib {
	return &DirsLib{}
}

func (d *DirsLib) GetConfig() string {
	config_dir := golangutils.ResolvePath(SystemUtils.HomeDir() + "/.config")
	golangutils.CreateDirectory(config_dir, true)
	return config_dir
}

func (d *DirsLib) GetOtherApps() string {
	return "other_apps_dir"
}
