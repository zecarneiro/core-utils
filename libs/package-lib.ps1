# ---------------------------------------------------------------------------- #
#                                    WINGET                                    #
# ---------------------------------------------------------------------------- #
# This function copied from the original: https://www.powershellgallery.com/packages/WingetTools/1.3.0
function install-winget {
    #Install the latest package from GitHub
    [cmdletbinding(SupportsShouldProcess)]
    [alias("iwg")]
    [OutputType("None")]
    [OutputType("Microsoft.Windows.Appx.PackageManager.Commands.AppxPackage")]
    Param(
        [Parameter(HelpMessage = "Display the AppxPackage after installation.")]
        [switch]$Passthru
    )

    if (!(commandexists "winget") -or (Read-Host "Winget is already installed, would you like to update it (y/N)") -eq "y") {
        infolog "Install Winget-CLI"
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
}

function install-winget-packages {
    evaladvanced "winget install --id=chrisant996.Clink --accept-source-agreements --accept-package-agreements"
}

# ---------------------------------------------------------------------------- #
#                                     SCOOP                                    #
# ---------------------------------------------------------------------------- #
function install-scoop {
    if (!(commandexists scoop)) {
        Write-Host "INFO: Install Scoop ..."
        Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
    }
}

function install-scoop-packages {
    evaladvanced "scoop install main/coreutils"
    evaladvanced "scoop install main/git"
    evaladvanced "scoop install main/vim"
    evaladvanced "scoop install main/nano"
    evaladvanced "scoop install main/curl"
    evaladvanced "scoop install main/grep"
    evaladvanced "scoop install main/sed"
    evaladvanced "scoop install main/which"
    evaladvanced "scoop install main/dos2unix"
    evaladvanced "scoop install main/7zip"
    evaladvanced "scoop install main/gsudo"
    evaladvanced "scoop install main/fzf"

    evaladvanced "scoop bucket add extras"
    evaladvanced "scoop install extras/psfzf"
    evaladvanced "scoop install extras/git-credential-manager"

    evaladvanced "scoop bucket add alkuzad_unxutils-separated https://github.com/alkuzad/unxutils-separated"
    evaladvanced "scoop install alkuzad_unxutils-separated/unxutils-xargs"

    # Markdown apps
    # evaladvanced "scoop install extras/marktext" # https://github.com/marktext/marktext
    evaladvanced "scoop install https://github.com/c3er/mdview/releases/latest/download/mdview.json" # https://github.com/c3er/mdview
}

function config-scoop-all {
    evaladvanced "gsudo config CacheMode auto"
    
    # Delete all system alias
    delalias "cp"
    delalias "cat"
    delalias "mkdir"
    delalias "ls"
    delalias "mv"
    delalias "ps"
    delalias "rm"
    delalias "rmdir"
    delalias "sleep"
    delalias "sort"
    delalias "tee"
    delalias "curl"
    delalias "grep"
    delalias "sed"

    # Docs
    titlelog "Integrate 7zip on context menu"
    log "Step 1: Open 7-Zip File Manager(As Admin) by typing 7-Zip in Start menu or Start screen and then pressing Enter key."
    log "Step 2: Next, navigate to Tools menu and then click Options to open Options dialog."
    log "Step 3: Here, under 7-Zip tab, make sure that Integrate 7-Zip to shell context menu option is selected. If not, please select the option and then click Apply button. You might need to reboot your PC or restart Windows Explorer to see 7-Zip in the context menu."
    pause
}

# ---------------------------------------------------------------------------- #
#                               POWERSHELL MODULE                              #
# ---------------------------------------------------------------------------- #
function install-powershellget {
    evaladvanced "sudopwsh Install-Module -Name PowerShellGet -Force"
}

function install-modules {
    evaladvanced "Install-Module -AllowClobber -Name scoop-completion -Scope CurrentUser -Force" # Project URL - https://github.com/Moeologist/scoop-completion"
    evaladvanced "Install-Module PSReadLine -Repository PSGallery -Scope CurrentUser -Force" # https://github.com/PowerShell/PSReadLine
}

# ---------------------------------------------------------------------------- #
#                                    OTHERS                                    #
# ---------------------------------------------------------------------------- #
function install-visual-c-runtimes {
    $batScript = "$BIN_DIR\Visual-C-Runtimes-All-in-One\install_all.bat"
    if ((__show_install_message_question "Visual-C-Runtimes") -eq "y") {
        infolog "Install Visual-C-Runtimes"
        evaladvanced "sudopwsh $batScript"
    }
}

# ------------------------------- DEPENDS SCOOP ------------------------------ #
function install-npm {
    infolog "Install NPM(Include Javascript by default)"
    evaladvanced "scoop bucket add main"
    evaladvanced "scoop install main/nodejs-lts"
}

# ------------------------ DEPENDS OF SCOOP AND PYTHON ----------------------- #
function install-pip-pipx {
    infolog "Install Pip/Pipx"
    warnlog "To continue, please install Python"
    __install_python
    evaladvanced "scoop install pipx"
    evaladvanced "pipx ensurepath --force"
    evaladvanced "pip install virtualenv"
}
