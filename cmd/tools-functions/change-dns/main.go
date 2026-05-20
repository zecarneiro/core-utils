package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/consolemenu"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"

	"main/cmd/tools-functions/change-dns/entities"
	"main/internal/libs"
)

var (
	consoleMenu *consolemenu.ConsoleMenu
	scriptFile  string
	servers     []entities.ServerIps
)

const (
	ipv4Type = "IPv4"
	ipv6Type = "IPv6"
)

func init() {
	console.EnableFeatures()
	if platform.IsWindows() && !system.IsAdmin() {
		logic.ProcessError(errors.New(system.NeedAdminAccessMsg))
	}
}

func loadVars() {
	consoleMenu = consolemenu.New()
	scriptFileExtension := logic.Ternary(platform.IsWindows(), ".ps1", ".sh")
	scriptFile = libs.GetScriptCmdPathByName(fmt.Sprintf(`change-dns%s`, scriptFileExtension), "tools-functions")
	servers = []entities.ServerIps{
		{
			Name: "google",
			Ipv4: entities.Ips{Primary: "8.8.8.8", Secundary: "8.8.4.4"},
			Ipv6: entities.Ips{Primary: "2001:4860:4860::8888", Secundary: "2001:4860:4860::8844"},
		},
		{
			Name: "quad9",
			Ipv4: entities.Ips{Primary: "9.9.9.9", Secundary: "149.112.112.112"},
			Ipv6: entities.Ips{Primary: "2620:fe::fe", Secundary: "2620:fe::9"},
		},
		{
			Name: "opendns",
			Ipv4: entities.Ips{Primary: "208.67.222.222", Secundary: "208.67.220.220"},
			Ipv6: entities.Ips{Primary: "2620:0:ccc::2", Secundary: "2620:0:ccd::2"},
		},
		{
			Name: "cloudflare",
			Ipv4: entities.Ips{Primary: "1.1.1.1", Secundary: "1.0.0.1"},
			Ipv6: entities.Ips{Primary: "2606:4700:4700::1111", Secundary: "2606:4700:4700::1001"},
		},
	}
}

func isValidIP(ip string) bool {
	return !str.IsEmpty(ip)
}

func isValidIPs(primary string, secundary string) bool {
	return isValidIP(primary) && isValidIP(secundary)
}

func runScript(operation string, scriptArgs ...string) {
	cmd := models.Command{
		Cmd:      fmt.Sprintf(`%s %s %s`, scriptFile, operation, slice.ArrayToString(scriptArgs)),
		UseShell: true,
		Verbose:  false,
	}
	logger.Error(exe.ExecRealTime(cmd))
}

func restartAdaperts() {
	logger.Separator()
	logger.Info("Restarting Adapters...")
	runScript("restartadapter")
}

func setServer(serverIp entities.ServerIps) {
	canRestart := false
	if isValidIPs(serverIp.Ipv4.Primary, serverIp.Ipv4.Secundary) {
		canRestart = true
		runScript("set", serverIp.Ipv4.Primary, serverIp.Ipv4.Secundary, ipv4Type)
	}
	if isValidIPs(serverIp.Ipv6.Primary, serverIp.Ipv6.Secundary) {
		canRestart = true
		runScript("set", serverIp.Ipv6.Primary, serverIp.Ipv6.Secundary, ipv6Type)
	}
	if canRestart {
		restartAdaperts()
	}
}

func setCustom(typeIp string) entities.Ips {
	message := "Insert IPs for %s - (PRESS ENTER TO CANCEL)"
	ip := entities.Ips{Primary: "", Secundary: ""}
	logger.Header(fmt.Sprintf(message, typeIp))
	fmt.Printf("Primary IP: ")
	fmt.Scanln(&ip.Primary)
	if isValidIP(ip.Primary) {
		fmt.Printf("Secundary IP: ")
		fmt.Scanln(&ip.Secundary)
	}
	return ip
}

func main() {
	loadVars()
	consoleMenu.Title = "Change DNS server settings"
	for _, dnsServer := range servers {
		consoleMenu.AddEntry(dnsServer.Name, func(entryMenu consolemenu.ConsoleMenuEntry) {
			setServer(dnsServer)
		})
	}
	consoleMenu.AddSeparator()
	consoleMenu.AddEntry("Custom", func(entryMenu consolemenu.ConsoleMenuEntry) {
		ipsv4 := setCustom(ipv4Type)
		ipsv6 := setCustom(ipv6Type)
		setServer(entities.ServerIps{Ipv4: ipsv4, Ipv6: ipsv6})
	})
	consoleMenu.AddEntry("Reset", func(entryMenu consolemenu.ConsoleMenuEntry) {
		runScript("reset")
		restartAdaperts()
	})
	consoleMenu.Start()
}
