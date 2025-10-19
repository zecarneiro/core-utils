from pathlib import Path

from libs.generic_libs import CONSOLE_UTILS
from vendor.pythonutils.file_utils import resolve_path
from vendor.pythonutils.entities.command_info import CommandInfo
from vendor.pythonutils.enums.shell_enum import EShell
from vendor.pythonutils.logger_utils import error_log
from vendor.pythonutils.system_utils import get_shell, get_os_name


def SHELL_PROFILE_SCRIPT():
    shell = get_shell()
    if shell == EShell.POWERSHELL or shell == EShell.CMD:
        cmd_res = CONSOLE_UTILS.exec(CommandInfo(command="$profile.CurrentUserAllHosts"))
        if len(cmd_res.stderr) > 0:
            error_log(cmd_res.stderr)
        else:
            print(cmd_res.stdout)
    else:
        profile_shell: str = f"{Path.home()}/"
        if shell == EShell.BASH:
            print(f"{profile_shell}.bashrc")
        elif shell == EShell.ZSH:
            print(f"{profile_shell}.zshrc")
        elif shell == EShell.FISH:
            print(f"{profile_shell}.config/fish/config.fish")
        else:
            print(EShell.UNKNOWN.value)

def CUSTOM_SHELL_PROFILE_SCRIPT():
    shell = get_shell()
    if shell == EShell.POWERSHELL or shell == EShell.CMD:
        print(resolve_path(f"{Path.home()}/.powershell-profile-custom.ps1"))
    else:
        shell_name = ""
        if shell == EShell.BASH:
            shell_name = "bash"
        elif shell == EShell.ZSH:
            shell_name = "zsh"
        elif shell == EShell.FISH:
            shell_name = "fish"

        if len(shell_name) > 0:
            print(f"{Path.home()}/.{shell_name}-profile-custom.sh")
        else:
            print(EShell.UNKNOWN.value)

def OS_NAME():
    print(get_os_name())

def ALIAS_SCRIPT():
    #declare MY_ALIAS = "$HOME/.bash_aliases"
    #$MY_ALIAS = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("$home\.powershell_aliases.ps1")
    print("AAAA")


