import sys
import json

from coreutils.entities.config import Config
from coreutils.libs.const_lib import CONSOLE_UTILS, SHELL_UTILS, SYSTEM_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils

__CONFIG__ = "config.json"
__CONFIG_FILE_NAME__ = "config.json"

from coreutils.libs.pythonutils.logger_utils import LoggerUtils

def get_args_str() -> str:
    return GenericUtils.list_to_str(sys.argv[1:])

def read_config() -> Config:
    config_file = FileUtils.resolve_path(f"{DirsLib.get_coreutils_config_dir()}/{__CONFIG_FILE_NAME__}")
    current_shell_name = SHELL_UTILS.current_shell.value
    try:
        with open(config_file, 'r') as f:
            data = json.load(f)
        config_data = Config.model_validate_json(data)
    except Exception as e:
        LoggerUtils.error_log(f"JSON validation: {e}")
        config_data: Config = Config(promptStyle={})
    if current_shell_name not in config_data.promptStyle:
        if SHELL_UTILS.is_powershell:
            config_data.promptStyle[current_shell_name] = 2
        else:
            config_data.promptStyle[current_shell_name] = 4
    FileUtils.write_file(config_file, GenericUtils.object_to_string(config_data))
    return config_data

def write_config(data: Config):
    config_file = FileUtils.resolve_path(f"{DirsLib.get_coreutils_config_dir()}/{__CONFIG_FILE_NAME__}")
    FileUtils.write_file(config_file, GenericUtils.object_to_string(data))

def get_all_shell_profiles_files() -> dict[EShell, str]:
    shells: dict[EShell, str] = {
        EShell.BASH: FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}/.bashrc"),
        EShell.ZSH: FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}/.zshrc"),
        EShell.FISH: FileUtils.resolve_path(f"{DirsLib.get_config()}/fish/config.fish"),
        EShell.KSH: FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}/.kshrc")
    }
    if SYSTEM_UTILS.is_windows:
        # This is only for POWERSHELL 7+ and none of actualy windows SO came with this version by default
        # resolve_path(f"{get_home_dir()}/Documents/PowerShell/Microsoft.PowerShell_profile.ps1")
        shells[EShell.POWERSHELL] = FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1")
    elif SYSTEM_UTILS.is_linux or SYSTEM_UTILS.is_macos:
        shells[EShell.POWERSHELL] = FileUtils.resolve_path(f"{DirsLib.get_config()}/powershell/Microsoft.PowerShell_profile.ps1")
    return shells

def set_file_permission_to_run(filepath: str):
    if FileUtils.is_file(filepath):
        if not SYSTEM_UTILS.is_windows:
            CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"chmod +x '{filepath}'"))
