import shutil
import sys
import argparse

from vendor.pythonutils.console_utils import ConsoleUtils
from vendor.pythonutils.entities.command_info import CommandInfo
from vendor.pythonutils.enums.shell_enum import EShell
from vendor.pythonutils.file_utils import file_exist, file_extension, filename_without_ext, delete_file, get_file_encoding
from vendor.pythonutils.logger_utils import get_bool_str_formated, info_log, error_log
from vendor.pythonutils.system_utils import is_windows, is_linux, get_shell


def gitresethardorigin():
    console_utils = ConsoleUtils()
    command_info = CommandInfo(command="git")

    # Get branch name
    command_info.args = ["branch", "--show-current"]
    response = console_utils.exec(command_info)
    if len(response.stderr) > 0:
        error_log(response.stderr)
    else:
        current_branch_name = response.stdout
        command_info.args = ["reset", "--hard", f"origin/{current_branch_name}"]
        console_utils.exec_real_time(command_info)
