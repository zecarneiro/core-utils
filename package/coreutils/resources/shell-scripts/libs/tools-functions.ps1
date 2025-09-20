# Author: José M. C. Noronha
# TODO: Implement kill-port
function runbashscriptWSL {
	param([string] $scriptOrigPath, [bool] $withSudo, [string] $distro, [string[]] $arguments = @())
	$scriptPath = (wintowslpath "$scriptOrigPath")
    if (![string]::IsNullOrEmpty($distro)) {
        $distro = "-d $distro"
    }
    $argumentsStr = ""
    foreach ($argument in $arguments) {
        $argumentsStr = "$argumentsStr '$argument'"
    }
	evaladvanced "wsl $distro -- chmod +x `"$scriptPath`""
	if ($withSudo) {
        $command = "sudo '$scriptPath' $argumentsStr"
	} else {
        $command = "source '$scriptPath' $argumentsStr"
	}
    evaladvanced "wsl $distro -- $command"
}
