import argparse
import os

from coreutils.libs.const_lib import CONSOLE_UTILS
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.systemfunctions import script_processor


def process_post_install_for_dir_function_file():
    script_proc: dict[EShell, list[str]]
    if CONSOLE_UTILS.shell_utils.is_shell([EShell.POWERSHELL, EShell.CMD]):
        script_processor(["-i", "-n", "ldir", "-c", "Get-ChildItem -Path \"$pwd\" -Directory -Force | ForEach-Object {$_.BaseName}"])
        script_processor(["-i", "-n", "countdirs", "-c", "(Get-ChildItem -Path \"$pwd\" -recurse | where-object { $_.PSIsContainer }).Count"])
        script_processor(["-i", "-n", "mkdir", "-c", f"New-Item -Path \"{CONST.POWERSHELL_ALL_ARGS_VAR_STR}\" -ItemType Directory -Force | Out-Null"])
        script_processor(["-i", "-n", "cpdir", "-c", f"param([string]$src, [string]$dest){CONST.EOF}Copy-Item \"$src\" -Destination \"$dest\" -Recurse -Force | Out-Null"])
    else:
        script_processor(["-i", "-n", "ldir", "-c", "find . -maxdepth 1 -type d -not -path ."])
        script_processor(["-i", "-n", "countdirs", "-c", "find . -type d | wc -l"])
        script_processor(["-i", "-n", "mkdir", "-c", f"mkdir -p {CONST.BASH_ALL_ARGS_VAR_STR}"])
        script_processor(["-i", "-n", "cpdir", "-c", f"cp -r {CONST.BASH_ALL_ARGS_VAR_STR}"])

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
    parser.add_argument("-d", "--directory", metavar="DIRECTORY_PATH",  type=str, required=True)
    args = parser.parse_args()
    print(LoggerUtils.get_bool_str_formated(FileUtils.is_dir(args.directory)))

def delete_directory():
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--directory", metavar="DIRECTORY_PATH", type=str, required=True)
    args = parser.parse_args()
    directory = args.directory
    if FileUtils.is_dir(directory):
        FileUtils.delete_file(directory)
        if FileUtils.is_dir(directory):
            LoggerUtils.error_log(f"Directory '{directory}' not deleted")

def find_dir():
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--name", metavar="DIR_NAME", type=str, required=True, help="Directory name to find")
    args = parser.parse_args()
    cwd = os.getcwd()
    if FileUtils.is_dir(cwd):
        dir_list = FileUtils.get_list_dirs_on_folder(cwd)
        for directory in dir_list:
            basedir = os.path.dirname(directory)
            if args.name in basedir:
                print(directory)
