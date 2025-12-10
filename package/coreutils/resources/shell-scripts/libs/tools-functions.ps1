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
function runExeFile {
    param([string] $exeFile)
    if ((fileexists "$exeFile") -or (commandexists "$exeFile")) {
        evaladvanced ""
    }
}
function runMsixFile {
    param([string] $exeFile)
    if ((fileexists "$exeFile")) {
        promptlog "Add-AppxPackage -Path `"$exeFile`""
        
    }
}
function runFilesByPath {
    param([string] $pathWithFiles, [string] $fileCache="")
    if ((directoryexists "$pathWithFiles")) {
        infolog "Sometimes it's necessary to run shell as admin."
        $exeExtArr = @("exe", "msixbundle", ".msi")
        $fileCacheData = ""
        if ((fileexists "$fileCache")) {
            $fileCacheData = (Get-Content "$fileCache")
        }
        foreach ($exeExt in $exeExtArr) {
            Get-ChildItem -Path "$pathWithFiles" -Filter "*.$exeExt" | ForEach-Object {
                $fileFull = $_.FullName
                if (($fileCacheData -notcontains "$fileFull")) {
                    if ("$exeExt" -eq "msixbundle") {
                        runMsixFile "$fileFull"
                    } elseif ("$exeExt" -eq "exe") {
                        runExeFile "$fileFull"
                    } elseif ("$exeExt" -eq "msi") {
                        runExeFile "$fileFull"
                    }
                    if (![string]::IsNullOrEmpty($fileCache)) {
                        writefile -file "$fileCache" -content "$fileFull" -append
                    }
                }
            }
        }
    }
    infolog "Execution of all files on '$pathWithFiles' it's done."
    pause
}
