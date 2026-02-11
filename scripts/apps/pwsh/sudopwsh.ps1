if ([string]::IsNullOrWhiteSpace($args)) {
    error-log "No args passed"
    exit
}
sudo powershell -Command $args