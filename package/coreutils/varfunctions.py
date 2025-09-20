from coreutils.libs.const_lib import SYSTEM_UTILS, SHELL_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.generic_lib import get_all_shell_profiles_files
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


def __print_data(data: str|None, status: bool, is_dir: bool):
    if data is None:
        LoggerUtils.warn_log(CONST.UNKNOWN)
    else:
        if status:
            print(data)
        else:
            LoggerUtils.error_log(f"Creating {"dir" if is_dir else "file"} fail: {data}")

def config_dir():
    directory = DirsLib.get_config()
    __print_data(directory, FileUtils.create_dir(directory), True)

def other_apps_dir():
    directory = DirsLib.get_user_opt()
    __print_data(directory, FileUtils.create_dir(directory), True)

def user_bin_dir():
    directory = DirsLib.get_user_bin()
    __print_data(directory, FileUtils.create_dir(directory), True)

def user_startup_dir():
    directory = DirsLib.get_user_startup()
    __print_data(directory, FileUtils.create_dir(directory), True)

def user_temp_dir():
    directory = DirsLib.get_user_temp()
    __print_data(directory, FileUtils.create_dir(directory), True)

def temp_dir():
    directory = DirsLib.get_temp()
    __print_data(directory, FileUtils.create_dir(directory), True)

def shell_profile_script():
    current_shell = SHELL_UTILS.current_shell
    if current_shell == EShell.UNKNOWN:
        LoggerUtils.warn_log(CONST.UNKNOWN)
    else:
        shells = get_all_shell_profiles_files()
        shell = shells.get(current_shell)
        if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
            shell = shells.get(EShell.POWERSHELL)
        status = True
        if not FileUtils.is_file(shell):
            status = FileUtils.touch(shell)
        __print_data(shell, status, False)

def shell_script_dir():
    directory = DirsLib.get_coreutils_shell_script()
    __print_data(directory, FileUtils.create_dir(directory), True)

def shell_alias_dir():
    directory = DirsLib.get_coreutils_shell_alias()
    __print_data(directory, FileUtils.create_dir(directory), True)

def shell_name():
    print(SHELL_UTILS.current_shell.value.upper())

def os_name():
    print(SYSTEM_UTILS.os_name)


