param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$RestoDosArgumentos
)

if ($RestoDosArgumentos.Count -gt 0) {
    $comando = $RestoDosArgumentos[0]
    $outros = $RestoDosArgumentos[1..($RestoDosArgumentos.Count - 1)]
    if ($comando -eq "install") {
        pip install --break-system-packages $outros
    } elseif ($comando -eq "uninstall") {
        pip uninstall --break-system-packages $outros
    } else {
        pip $RestoDosArgumentos
    }
} else {
    pip
}