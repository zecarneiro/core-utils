import sys

from coreutils.entities.config import Config
from coreutils.libs.const_lib import CONSOLE_UTILS, SHELL_UTILS, SYSTEM_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils

__CONFIG__ = "config.json"
__CONFIG_FILE_NAME__ = "config.json"


def get_args_str() -> str:
    return GenericUtils.list_to_str(sys.argv[1:])


def read_config() -> Config:
    config_file = FileUtils.resolve_path(
        f"{DirsLib.get_coreutils_config_dir()}/{__CONFIG_FILE_NAME__}"
    )
    current_shell_name = SHELL_UTILS.current_shell.value
    config_data = FileUtils.load_json_file(config_file, Config)
    if config_data is None:
        config_data = Config(promptStyle={})
    if current_shell_name not in config_data.promptStyle:
        if SHELL_UTILS.is_powershell:
            config_data.promptStyle[current_shell_name] = 2
        else:
            config_data.promptStyle[current_shell_name] = 4
    FileUtils.write_file(config_file, GenericUtils.object_to_string(config_data))  # type: ignore[reportUnknownMemberType]
    return config_data


def write_config(data: Config):
    config_file = FileUtils.resolve_path(
        f"{DirsLib.get_coreutils_config_dir()}/{__CONFIG_FILE_NAME__}"
    )
    FileUtils.write_file(config_file, GenericUtils.object_to_string(data))  # type: ignore[reportUnknownMemberType]


def set_file_permission_to_run(filepath: str):
    if FileUtils.is_file(filepath):
        if not SYSTEM_UTILS.is_windows:
            CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"chmod +x '{filepath}'"))
