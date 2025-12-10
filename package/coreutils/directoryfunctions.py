import argparse
import os

from coreutils.libs.const_lib import SHELL_UTILS, CONSOLE_UTILS
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


def delete_empty_dirs():
    for dirpath, _, _ in os.walk(os.getcwd(), topdown=False):
        try:
            if not os.listdir(dirpath):
                os.rmdir(dirpath)
                LoggerUtils.ok_log(f"Deleted: {dirpath}")
        except (FileNotFoundError, PermissionError) as e:
            LoggerUtils.warn_log(f"Skipped {dirpath}: {e}")

def directory_exists():
    parser = argparse.ArgumentParser()
    parser.add_argument("directory",  type=str)
    args = parser.parse_args()
    directory = args.directory
    if directory is None or directory == ".":
        directory = ""
    print(LoggerUtils.get_bool_str_formated(FileUtils.is_dir(directory)))

def delete_directory():
    parser = argparse.ArgumentParser()
    parser.add_argument("directory", type=str)
    args = parser.parse_args()
    directory = args.directory
    if FileUtils.is_dir(directory):
        FileUtils.delete_file(directory)
        if FileUtils.is_dir(directory):
            LoggerUtils.error_log(f"Directory '{directory}' not deleted")
    else:
        LoggerUtils.error_log(f"Not found directory: {directory}")

def find_dir():
    parser = argparse.ArgumentParser()
    parser.add_argument("name", type=str, help="Directory name to find")
    args = parser.parse_args()
    cwd = os.getcwd()
    if args.name is not None and len(args.name) > 0:
        if FileUtils.is_dir(cwd):
            dir_list = FileUtils.get_list_dirs_on_folder(cwd)
            for directory in dir_list:
                basedir = os.path.dirname(directory)
                if args.name in basedir:
                    print(directory)
    else:
        LoggerUtils.error_log(f"Invalid given name: {args.name}")

def lh_dir():
    if SHELL_UTILS.is_powershell:
        cmd = "Get-ChildItem -Path \"$pwd\" -Directory -Force | ForEach-Object {$_.BaseName}"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, verbose=True, shell=EShell.POWERSHELL, use_shell=True))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command="find . -maxdepth 1 -type d -not -path .", verbose=True))

def count_dirs():
    if SHELL_UTILS.is_powershell:
        cmd = "(Get-ChildItem -Path \"$pwd\" -recurse | where-object { $_.PSIsContainer }).Count"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL, verbose=True, use_shell=True))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command="find . -type d | wc -l", verbose=True))

def mkdirc():
    parser = argparse.ArgumentParser()
    parser.add_argument("directory", type=str, help="Directory to create. Create if exist and recursive")
    args = parser.parse_args()
    directory = args.directory
    if SHELL_UTILS.is_powershell:
        cmd = f"New-Item -Path \"{directory}\" -ItemType Directory -Force | Out-Null"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL, verbose=True, use_shell=True))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command=f"mkdir -p {directory}", verbose=True))

def cp_dir():
    parser = argparse.ArgumentParser()
    parser.add_argument("-s", "--source", metavar="SOURCE_DIR", type=str, required=True)
    parser.add_argument("-d", "--dest", metavar="DEST_DIR", type=str, required=True)
    args = parser.parse_args()
    source = args.source
    dest = args.dest
    if SHELL_UTILS.is_powershell:
        cmd = f"Copy-Item \"{source}\" -Destination \"{dest}\" -Recurse -Force | Out-Null"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL, verbose=True, use_shell=True))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command=f"cp -r \"{source}\" \"{dest}\"", verbose=True))