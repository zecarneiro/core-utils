package processors

import (
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/netc"
	"golangutils/pkg/system"
	"strings"
)

type Downloader struct {
	force bool
}

func (dp *Downloader) getExtensionVsixPath() string {
	dir := file.JoinPath(system.TempDir(), "vscode-config-download")
	logic.ProcessError(file.CreateDirectory(dir, true))
	return dir
}

func (dp *Downloader) getExtensionVsixFile(extensionId string) string {
	path := dp.getExtensionVsixPath()
	if len(path) == 0 {
		return path
	}
	extensionDest := path + "/" + extensionId + ".vsix"
	return file.ResolvePath(extensionDest)
}

func (dp *Downloader) getUrl(extensionId string) string {
	delimiter := "."
	url := "https://{publisher}.gallery.vsassets.io/_apis/public/gallery/publisher/{publisher}/extension/{package}/latest/assetbyname/Microsoft.VisualStudio.Services.VSIXPackage"
	// Split the string into substrings using the delimiter
	publisherPackage := strings.Split(extensionId, delimiter)
	publisherExtension := publisherPackage[0]
	packageExtension := publisherPackage[1]
	newUrl := strings.Replace(url, "{publisher}", publisherExtension, -1)
	newUrl = strings.Replace(newUrl, "{package}", packageExtension, -1)
	return newUrl
}

func (dp *Downloader) download(extensionId string) bool {
	processor := &Downloader{}
	status := false
	filePath := processor.getExtensionVsixFile(extensionId)
	if dp.force {
		logger.Error(file.DeleteFile(filePath))
	}
	if len(filePath) > 0 && !file.IsFile(filePath) {
		if status, err := netc.HasInternet(); !status {
			logger.ErrorStr("Not detect internet.")
			logger.Error(err)
			console.WaitForAnyKeyPressed("Please, connect to internet(PRESS ANY KEY TO CONTINUE)")
		}
		err := netc.Download(processor.getUrl(extensionId), filePath)
		if err != nil {
			logger.Error(err)
		} else {
			status = true
		}
	}
	return status
}

func (dp *Downloader) downloadList(extensionIdList []string) {
	for _, id := range extensionIdList {
		dp.download(id)
	}
}
