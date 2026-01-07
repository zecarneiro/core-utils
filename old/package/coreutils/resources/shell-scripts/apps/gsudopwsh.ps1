if ((text-is-empty "$args") -eq "true") {
    error-log "No args passed"
    exit
}
gsudo powershell -Command $args