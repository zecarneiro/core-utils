param(
    [string] $expression
)

if ([string]::IsNullOrEmpty($expression)) {
    Get-ChildItem env:*
}
else {
    if (!$expression.Contains("=")) {
        Write-Output "Environment variable $expression not defined"
    }
    else {
        $expressionArr = $expression.Split("=")
        $name = $expressionArr[0]
        $value = ""
        if ($expressionArr.Length -gt 1) {
            $value = $expressionArr[1]
        }
        set-item -force -path "env:$name" -value $value
    }
}