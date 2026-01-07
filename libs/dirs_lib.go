package libs

import "golangutils"

type DirsLib struct{}

func NewDirsLib() *DirsLib {
	return &DirsLib{}
}

func (d *DirsLib) GetUserConfig() string {
	config_dir := golangutils.ResolvePath(SystemUtils.HomeDir() + "/.config")
	golangutils.CreateDirectory(config_dir, true)
	return config_dir
}

func (d *DirsLib) GetUserLocal() string {
	config_dir := golangutils.ResolvePath(SystemUtils.HomeDir() + "/.local")
	golangutils.CreateDirectory(config_dir, true)
	return config_dir
}

func (d *DirsLib) GetUserOpt() string {
	opt_dir := golangutils.ResolvePath(d.GetUserLocal() + "/opt")
	golangutils.CreateDirectory(opt_dir, true)
	return opt_dir
}
