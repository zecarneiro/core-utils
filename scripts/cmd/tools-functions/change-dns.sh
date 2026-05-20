#!/usr/bin/env bash

declare CONNECTION_NAMES=($(nmcli device status | grep -i -e ethernet -i -e wireguard -i -e wifi | awk '{ print $4}'))
declare OPERATION_ARG="$1"
declare IP1="$2"
declare IP2="$3"
declare IP_TYPE="$4"

restartAdapter() {
    sudo systemctl restart NetworkManager.service
    ok-log "Restart Systemd and Network services"
}

setDns() {
    local ip="${IP1},${IP2}"
    if [[ -z "${IP1}" ]]||[[ -z "${IP2}" ]]; then
        error-log "Invalid given DNS IPs"
        exit 1
    fi    
    if [[ "${IP_TYPE}" != "IPv4" ]]&&[[ "${IP_TYPE}" != "IPv6" ]]; then
        error-log "Invalid given type of IPs(IPv4 or IPv6)"
        exit 1
    fi
    for connectionName in "${CONNECTION_NAMES[@]}"; do
        if [[ -n "${connectionName}" ]]&&[[ "${connectionName}" != "--" ]]; then
            if [[ "${IP_TYPE}" == "IPv4" ]]; then
                nmcli connection modify "${connectionName}" ipv4.dns "$ip" ipv4.ignore-auto-dns yes
            fi
            if [[ "${IP_TYPE}" == "IPv6" ]]; then
                nmcli connection modify "${connectionName}" ipv6.dns "$ip" ipv6.ignore-auto-dns yes
            fi
            ok-log "$IP_TYPE DNS settings was changed for ${connectionName} to $IP1 and $IP2"
        fi
    done
}

resetDns() {
    for connectionName in "${CONNECTION_NAMES[@]}"; do
        if [[ -n "${connectionName}" ]]&&[[ "${connectionName}" != "--" ]]; then
            nmcli connection modify "${connectionName}" ipv4.dns "" ipv4.ignore-auto-dns false
            nmcli connection modify "${connectionName}" ipv6.dns "" ipv6.ignore-auto-dns false
            ok-log "DNS settings was reseted for ${connectionName}"
        fi
    done
}

case "${OPERATION_ARG}" in
    reset) resetDns ;;
    restartadapter) restartAdapter ;;
    set) setDns ;;
    *) warn-log "Invalid given operation" ;;
esac

