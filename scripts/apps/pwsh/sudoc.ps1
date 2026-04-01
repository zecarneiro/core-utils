if ([string]::IsNullOrWhiteSpace($args))
{
    error-log "No args passed"
    exit 1
}
prompt-log "$args"
sudopwsh "$args"
