#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2164

function cutadvanced {
    local delimiter data direction # direction = L/R
    CUT_PARSED_ARGUMENTS=$(getopt --longoptions data:,direction:,delimiter: -o "" -- "$@")
    CUT_VALID_ARGUMENTS=$?
    if [ "$CUT_VALID_ARGUMENTS" != "0" ]; then
        errorlog "Invalid arguments on Cut command"
        return 1
    fi
    eval set -- "$CUT_PARSED_ARGUMENTS"
    while :; do
        case "$1" in
            --direction)
                direction=$2
                shift 2
            ;;
            --delimiter)
                delimiter=$2
                shift 2
            ;;
            --data)
                data=$2
                shift 2
            ;;
            --help)
                log "cutadvanced --data DATA --delimiter DELIMITIER --direction [L|R]"
                return
            ;;
            --)
                shift
                break
            ;;
        esac
    done
    if [[ -n "${data}" ]]; then
        if [[ "${direction}" == "R" ]]; then
            data="$(echo "$data" | awk -F "$delimiter" '{ print $2 }')"
            elif [[ "${direction}" == "L" ]]; then
            data="$(echo "$data" | awk -F "$delimiter" '{ print $1 }')"
        fi
        echo "$data"
    fi
}
function extract {
    local file="$1"
    local destination="$2"

    if [[ -z "${destination}" ]]; then
        destination="$PWD"
    fi
    local currentDirectory="$destination"
    mkdir "$destination"
    cd "$destination"
    if [ -f "$file" ] ; then
        case $file in
            *.tar.bz2)   tar xvjf "$file"    ;;
            *.tar.gz)    tar xvzf "$file"    ;;
            *.bz2)       bunzip2 "$file"     ;;
            *.rar)       rar x "$file"       ;;
            *.gz)        gunzip "$file"      ;;
            *.tar)       tar xvf "$file"     ;;
            *.tbz2)      tar xvjf "$file"    ;;
            *.tgz)       tar xvzf "$file"    ;;
            *.zip)       unzip "$file"       ;;
            *.Z)         uncompress "$file"  ;;
            *.7z)        7z x "$file"        ;;
            *)           infolog "don't know how to extract '$file'..." ;;
        esac
    else
        errorlog "'$file' is not a valid file!"
    fi
    cd "$currentDirectory"
}
function openurl {
    local url="$1"
    if [[ -n "${url}" ]]; then
        xdg-open "$url"
    fi
}
function hasinternet {
    if ! ping -c 1 8.8.8.8 -q &>/dev/null || ! ping -c 1 8.8.4.4 -q &>/dev/null || ! ping -c 1 time.google.com -q &>/dev/null; then
        echo false
    else
        echo true
    fi
}
function mypubip {
    # Dumps a list of all IP addresses for every device
    # /sbin/ifconfig |grep -B1 "inet addr" |awk '{ if ( $1 == "inet" ) { print $2 } else if ( $2 == "Link" ) { printf "%s:" ,$1 } }' |awk -F: '{ print $1 ": " $3 }';

    ### Old commands
    # Internal IP Lookup
    #echo -n "Internal IP: " ; /sbin/ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}'

    # External IP Lookup
    #echo -n "External IP: " ; wget http://smart-ip.net/myip -O - -q

    # Internal IP Lookup.
    if [ -e /sbin/ip ];
    then
        echo -n "Internal IP: " ; /sbin/ip addr show wlan0 | grep "inet " | awk -F: '{print $1}' | awk '{print $2}'
    else
        echo -n "Internal IP: " ; /sbin/ifconfig wlan0 | grep "inet " | awk -F: '{print $1} |' | awk '{print $2}'
    fi

    # External IP Lookup
    echo -n "External IP: " ; curl -s ifconfig.me
}
function download {
    local url
    local file
    while [ "$#" -ne 0 ] ; do
        case "${1}" in
            --url) url="$2"; shift 2 ;;
            --file) file="$2"; shift 2 ;;
						--help) log "download -url URL -file FILE"; return ;;
            *) shift ;;
        esac
    done
    if [ "$(hasinternet)" == false ]; then
        errorlog "No Internet connection available"
    else
        if [ "$(commandexists "wget")" == false ]; then
            evaladvanced "sudo apt install wget -y"
        fi
        infolog "Downloading from URL: $url"
        wget -O "$file" "$url" -q --show-progress
    fi
}
function exitwithmsg {
    local message="$1"
    local code=$2
    if [[ "--help" == "${message}" ]]||[[ "-h" == "${message}" ]]; then
        log "exitwithmsg MSG CODE(Default = 0)"
        return
    fi
    if [[ -n "${message}" ]]; then
        infolog "$message"
    fi
    if [[ -z "${code}" ]]; then
        code=0
    fi
    exit $code
}
alias now='date'
alias lhiden='ls -d .* --color=auto'
function runlineascommand {
    local file="$1"
    local headerKey="$2"
	local prefix_sufix_key="######"
	local allKey="$prefix_sufix_key ALL $prefix_sufix_key"
	local headerKey="$prefix_sufix_key $headerKey $prefix_sufix_key"
	local canRun=false
    while read line; do
        if [[ "${line}" == "${prefix_sufix_key}"* ]]; then
            if [[ "${line}" == "${allKey}"* ]]||[[ "${line}" == "${headerKey}"* ]]; then
                canRun=true
            else
                canRun=false
            fi
        fi
        if [[ "${canRun}" == "true" ]]&&[[ -n "${line}" ]]; then            
            if [[ "${line}" != "${prefix_sufix_key}"* ]]; then
                evaladvanced "$line" 
            fi
        fi
    done <"$file"
}
function waituntil {
    local seconds="$1"
    if [[ -z "${seconds}" ]]; then
        seconds=0
    fi
    echo ""
    read -t $seconds -n 1 -s -r -p "Waiting for ${seconds} seconds, press any key to continue ..."; echo " "
}
function dos2unixrec {
    local ext="$1"
    find . -type f -name "*.${ext}" -print0 | xargs -0 dos2unix
}
function unix2dosrec {
    local ext="$1"
    find . -type f -name "*.${ext}" -print0 | xargs -0 unix2dos
}
function settrycatch {
    local tryContent="$1"
    local catchContent="$2"
    local file="$3"
    writefile "$file" "($tryContent)||($catchContent)" -append
}
function kill-port {
    local port="$1"
    sudo kill -9 $(sudo lsof -t -i :$port)
}

# Feature for linux only
alias sha1='openssl sha1'
alias md5='openssl md5'
alias sha256='openssl sha256'
function pause {
    echo -n "Press [ENTER] to continue...: "
    read var_name
}
alias restart-pipewire="systemctl --user restart pipewire.service"
function nautilus-install-script-context-menu {
    local scriptName="$1"; shift
    local scriptCommands=( "$@" )
    local scriptsPath="$HOME/.local/share/nautilus/scripts"
    local scriptInstall="${scriptsPath}/${scriptName}"
    evaladvanced "mkdir -p '$scriptsPath'"
    nautilus-uninstall-script-context-menu "$scriptName"

    echo "#!/usr/bin/env bash" > "$scriptInstall"
    for scriptCommand in "${scriptCommands}"; do
        infolog "Insert into $scriptName the command: $scriptCommand"
        echo "$scriptCommand" | tee -a "$scriptInstall" >/dev/null
    done
    evaladvanced "chmod +x '$scriptInstall'"
    restartexplorer
}
function nautilus-uninstall-script-context-menu {
    local scriptName="$1"
    local scriptsPath="~/.local/share/nautilus/scripts/$scriptName"
    local scriptInstall="${scriptsPath}/$(basename "$script")"
    if [[ -f "$scriptsPath" ]]; then
        evaladvanced "rm -rf '$scriptsPath'"
        restartexplorer
    fi
}
function changedefaultjdk {
    local java_default_script_name="/etc/profile.d/jdk-default.sh"
    evaladvanced "update-java-alternatives --list"
    read -p "Insert Path of java(JAVA_HOME): " javaHome
    evaladvanced "echo \"JAVA_HOME_DEFAULT='${javaHome}'\" | sudo tee ${java_default_script_name}"
    evaladvanced "echo \"export JAVA_HOME=${javaHome}\" | sudo tee -a ${java_default_script_name}"
    evaladvanced "echo \"export PATH=\\$PATH:\\${JAVA_HOME_DEFAULT}/bin\" | sudo tee -a ${java_default_script_name}"
    evaladvanced "source ${java_default_script_name}"
}
alias chmod-777="chmod -R 777"
