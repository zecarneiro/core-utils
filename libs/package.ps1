function __set_user_bin_dir {
    $user_bin_dir = "$home\.local\bin"
    $pathKey = "Path"
    $pathEnvArr = ([Environment]::GetEnvironmentVariable($pathKey, [System.EnvironmentVariableTarget]::User) -split ';')

    __print "Set User bin dir"
    __eval "New-Item -ItemType Directory -Path `"$user_bin_dir`" -Force | Out-Null"
    if (!("$user_bin_dir" -in $pathEnvArr)) {
        $pathEnvArr += "$user_bin_dir"
        [Environment]::SetEnvironmentVariable($pathKey, ($pathEnvArr -join ";"), [System.EnvironmentVariableTarget]::User)
        Write-Host "[INFO] Please, Restart the Terminal to change take effect!"
    }
}

# This function copied from the original: https://www.powershellgallery.com/packages/WingetTools/1.3.0
function __install_winget {
    #Install the latest package from GitHub
    [cmdletbinding(SupportsShouldProcess)]
    [alias("iwg")]
    [OutputType("None")]
    [OutputType("Microsoft.Windows.Appx.PackageManager.Commands.AppxPackage")]
    Param(
        [Parameter(HelpMessage = "Display the AppxPackage after installation.")]
        [switch]$Passthru
    )
    __print "Install Winget-CLI"
    Write-Verbose "[$((Get-Date).TimeofDay)] Starting $($myinvocation.mycommand)"

    if ($PSVersionTable.PSVersion.Major -eq 7) {
        Write-Warning "This command does not work in PowerShell 7. You must install in Windows PowerShell."
        return
    }

    #test for requirement
    $Requirement = Get-AppPackage "Microsoft.DesktopAppInstaller"
    if (-Not $requirement) {
        Write-Verbose "Installing Desktop App Installer requirement"
        Try {
            Add-AppxPackage -Path "https://aka.ms/Microsoft.VCLibs.x64.14.00.Desktop.appx" -erroraction Stop
        }
        Catch {
            Throw $_
        }
    }

    $uri = "https://api.github.com/repos/microsoft/winget-cli/releases"

    Try {
        Write-Verbose "[$((Get-Date).TimeofDay)] Getting information from $uri"
        $get = Invoke-RestMethod -uri $uri -Method Get -ErrorAction stop

        Write-Verbose "[$((Get-Date).TimeofDay)] getting latest release"
        #$data = $get | Select-Object -first 1
        $data = $get[0].assets | Where-Object name -Match 'msixbundle'

        $appx = $data.browser_download_url
        #$data.assets[0].browser_download_url
        Write-Verbose "[$((Get-Date).TimeofDay)] $appx"
        If ($pscmdlet.ShouldProcess($appx, "Downloading asset")) {
            $file = Join-Path -path $env:temp -ChildPath $data.name

            Write-Verbose "[$((Get-Date).TimeofDay)] Saving to $file"
            Invoke-WebRequest -Uri $appx -UseBasicParsing -DisableKeepAlive -OutFile $file

            Write-Verbose "[$((Get-Date).TimeofDay)] Adding Appx Package"
            Add-AppxPackage -Path $file -ErrorAction Stop

            if ($passthru) {
                Get-AppxPackage microsoft.desktopAppInstaller
            }
        }
    } #Try
    Catch {
        Write-Verbose "[$((Get-Date).TimeofDay)] There was an error."
        Throw $_
    }
    Write-Verbose "[$((Get-Date).TimeofDay)] Ending $($myinvocation.mycommand)"
}

function __install_scoop_and_packages {
    __print "Install Scoop"
    __eval "Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression"
    __eval "Install-Module -AllowClobber -Name scoop-completion -Scope CurrentUser -Force" # Project URL - https://github.com/Moeologist/scoop-completion"

    __eval ". `"$__SHELL_FILE__`"" $true
    . "$__SHELL_FILE__"

    __eval "scoop bucket add extras"
    __eval "scoop install main/gsudo"
    __eval "scoop install main/topgrade"
    __eval "scoop install main/fzf"
    __eval "scoop install main/dos2unix"
    __eval "scoop install extras/psfzf"
    __eval "scoop install extras/psreadline" # https://github.com/PowerShell/PSReadLine

    __print "Install Python and PIP/PIPX"
    __eval "scoop install main/python"
    __eval "scoop install pipx"
    __eval "pipx ensurepath"
    __eval "pip install virtualenv"
}

function __install_visual_c_runtimes {
    $batScript = "$__BIN_DIR__\Visual-C-Runtimes-All-in-One\install_all.bat"
    __print "Install Visual-C-Runtimes"
    evaladvanced "gsudo $batScript"
}

function __install_features_for_wsl {
    if ((__show_install_message_question "Features for WSL") -eq "y") {
        infolog "Enable Virtual Machine Platform feature"
        evaladvanced "sudopwsh dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart"
        infolog "Enable WSL feature"
        evaladvanced "sudopwsh dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart"
    }
}