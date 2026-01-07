param(
    [string] $file,
    [string] $section
)
$prefix_sufix_key="######"
$allKey="$prefix_sufix_key ALL $prefix_sufix_key"
$section="$prefix_sufix_key $section $prefix_sufix_key"
$canRun=$false
foreach ($line in Get-Content "$file") {
    if(($line -like "${prefix_sufix_key}*")) {
        if (($line -like "${allKey}*") -or ($line -like "${section}*")) {
            $canRun=$true
        } else {
            $canRun=$false
        }
    }
    if ($canRun -and ![string]::IsNullOrEmpty($line)) {
        if (!($line -like "${prefix_sufix_key}*")) {
            evalc"$line"
        }
    }
}