package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/netc"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	urlArg           string
	outputFile       string
	secondVersionArg bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a file from URL",
	}
	cobralib.CobraCmd.Flags().StringVarP(&urlArg, "url", "u", "", "URL to download (required)")
	cobralib.CobraCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Output file path (required)")
	cobralib.CobraCmd.Flags().BoolVarP(&secondVersionArg, "version-2", "2", false, "Will run the system app version(Windows: Invoke-WebRequest/Linux: wget)")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("url"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("output-file"))
	cobralib.WithRun(process)
}

func process() {
	if str.IsEmpty(urlArg) {
		logic.ProcessError(fmt.Errorf("invalid given url"))
	}
	if str.IsEmpty(outputFile) {
		logic.ProcessError(fmt.Errorf("invalid given output file"))
	}
	if !platform.IsWindows() && !platform.IsLinux() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	status, err := netc.HasInternet()
	logic.ProcessError(err)
	if !status {
		logic.ProcessError(fmt.Errorf("no Internet connection available"))
	}
	if !secondVersionArg {
		logic.ProcessError(netc.Download(urlArg, outputFile))
	} else {
		cmd := models.Command{
			Verbose:  true,
			UseShell: true,
		}
		if platform.IsWindows() {
			cmd.Cmd = fmt.Sprintf("Invoke-WebRequest -Uri \"%s\" -OutFile \"%s\"", urlArg, outputFile)
		} else {
			cmd.Cmd = fmt.Sprintf("wget -O \"%s\" \"%s\" -q --show-progress", outputFile, urlArg)
		}
		logic.ProcessError(exe.ExecRealTime(cmd))
	}
}

func main() {
	cobralib.Run()
}
