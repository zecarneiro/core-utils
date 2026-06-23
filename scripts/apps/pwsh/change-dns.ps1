# Based on https://github.com/HarmVeenstra/Powershellisfun
# Author: José Noronha

# Need Admin to run
if ($(is-admin) -eq "false") {
    sudopwsh "$PSCommandPath"
    exit 0
}

# Global variables
$adapters = Get-NetAdapter | Where-Object Status -eq 'Up'
$computerName = $env:COMPUTERNAME
$dns_ipv4 = @{
    "google"     = @("8.8.8.8,8.8.4.4")
    "quad9"      = @("9.9.9.9,149.112.112.112")
    "opendns"    = @("208.67.222.222,208.67.220.220")
    "cloudflare" = @("1.1.1.1,1.0.0.1")
}
$dns_ipv6 = @{
    "google"     = @("2001:4860:4860::8888,2001:4860:4860::8844")
    "quad9"      = @("2620:fe::fe::8888,2620:fe::9")
    "opendns"    = @("2620:0:ccc::2,2620:0:ccd::2")
    "cloudflare" = @("2606:4700:4700::1111,2606:4700:4700::1001")
}

function restartAdapters {
    foreach ($adapter in $adapters) {
        try {
            $adapterName = $adapter.Name
            sudopwsh Restart-NetAdapter -Name "$adapterName"
            ok-log "Restart NetAdapter with name $adapterName, done."
        } catch {
            error-log "$($_.Exception.Message)"
        }
    }
}

function setDns {
    param([string]$ip, [string]$ip_type)
    if ($ip_type -ne "IPv4" -and $ip_type -ne "IPv6") {
        error-log "Invalid given type of IPs(IPv4 or IPv6)"
        exit 1
    }
    foreach ($adapter in $adapters) {
        if (-not [string]::IsNullOrEmpty($ip)) {
            $adapterName = $adapter.Name
            try {
                sudopwsh Set-DNSClientServerAddress -InterfaceIndex $adapter.ifIndex -ServerAddresses "$ip" -ErrorAction Stop
                ok-log "$ip_type DNS settings was changed for '$adapterName' to '$ip'"
            } catch {
                error-log "$($_.Exception.Message)"
            }
        }
    }
}

function setDefault {
    header-log "Change IPv4 and IPv6 DNS server settings"
    $providers = $dns_ipv4.Keys | Sort-Object
    while ($true) {
        for ($i = 0; $i -lt $providers.Count; $i++) {
            Write-Host "$($i + 1). $($providers[$i])"
        }
        Write-Host "$($providers.Count + 1). Exit"
        $selection = Read-Host "Select a provider"
        if ($selection -eq ($providers.Count + 1)) {
            break
        }
        $index = [int]$selection - 1
        if ($index -ge 0 -and $index -lt $providers.Count) {
            $provider = $providers[$index]
            $ipv4 = $dns_ipv4[$provider]
            $ipv6 = $dns_ipv6[$provider]
            Write-Host "Server selected: $provider"
            separator-log
            setDns "$ipv4" "IPv4"
            setDns "$ipv6" "IPv6"
            restartAdapters
            break
        }
        else {
            warn-log "❌ Invalid option. Try again."
        }
    }
}

function resetDns {
    foreach ($adapter in $adapters) {
        $adapterName = $adapter.Name
        if (Get-DnsClientServerAddress | Where-Object InterfaceIndex -eq $adapter.InterfaceIndex) {
            try {
                sudopwsh Set-DNSClientServerAddress -InterfaceIndex $adapter.ifIndex -ResetServerAddresses -ErrorAction Stop
                ok-log "DNS settings was reseted for $adapterName on $computerName"
            } catch {
                error-log "$($_.Exception.Message)"
            }
        }
    }
    restartAdapters
}

while ($true) {
    title-log "Change DNS server settings"
    Write-Host "1. Set Default Server"
    Write-Host "2. Reset"
    Write-Host "3. Exit"
    $choice = Read-Host "Insert an option"
    switch ($choice) {
        "1" { setDefault }
        "2" { resetDns }
        "3" { exit }
        default {warn-log "Invalid option inserted" }
    }
}
