# Based on https://github.com/HarmVeenstra/Powershellisfun
# Author: José Noronha
param(
    [Parameter(Mandatory=$true)]
    [string]$OPERATION_ARG,
    [string] $IP1,
    [string] $IP2,
    [ValidateSet("IPv4","IPv6")]
    [string]$IP_TYPE
)

# Global variables
$adapters = Get-NetAdapter | Where-Object Status -eq 'Up'
$computerName = $env:COMPUTERNAME

function restartAdapter {
    foreach ($adapter in $adapters) {
        try {
            $adapterName = $adapter.Name
            Restart-NetAdapter -Name "$adapterName"
            ok-log "Restart NetAdapter with name $adapterName, done."
        } catch {
            error-log "Changing $adapterName on $computerName"
        }
    }
}

function setDns {
    if ([string]::IsNullOrEmpty($IP1) -or [string]::IsNullOrEmpty($IP2)) {
        error-log "Invalid given DNS IPs"
        exit 1
    }
    #Loop through all adapters and configure $IP1 and $IP2 for all adapters which have a DNS Server setting
    foreach ($adapter in $adapters) {
        if (Get-DnsClientServerAddress | Where-Object InterfaceIndex -eq $adapter.InterfaceIndex) {
            $adapterName = $adapter.Name
            try {
                Set-DNSClientServerAddress -InterfaceIndex $adapter.ifIndex -ServerAddresses @($IP1, $IP2) -ErrorAction Stop
                ok-log "$IP_TYPE DNS settings was changed for $adapterName to $IP1 and $IP2 on $computerName"
            } catch {
                error-log "$IP_TYPE DNS: Changing $adapterName on $computerName"
            }
        }
    }
}

function resetDns {
    foreach ($adapter in $adapters) {
        $adapterName = $adapter.Name
        if (Get-DnsClientServerAddress | Where-Object InterfaceIndex -eq $adapter.InterfaceIndex) {
            try {
                Set-DNSClientServerAddress -InterfaceIndex $adapter.ifIndex -ResetServerAddresses -ErrorAction Stop
                ok-log "DNS settings was reseted for $adapterName on $computerName"
            } catch {
                error-log "Reset settings $adapterName on $computerName"
            }
        }
    }
}

switch ($OPERATION_ARG) {
    reset { resetDns }
    restartadapter { restartAdapter }
    set { setDns }
    Default { warn-log "Invalid given operation" }
}
