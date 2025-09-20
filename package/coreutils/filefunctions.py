import argparse
import os
import shutil

from coreutils.libs.const_lib import SYSTEM_UTILS, CONSOLE_UTILS, SHELL_UTILS
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.entities.write_file_options import WriteFileOptions
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


def file_dependencies_apps() -> list[str]:
    apps: list[str] = ["markdown_viewer"]
    if SYSTEM_UTILS.is_linux:
        apps.append("eog")
    return apps

def file_exists():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file = args.file
    if file is None or file == ".":
        file = ""
    print(LoggerUtils.get_bool_str_formated(FileUtils.is_file(file)))

def file_extension():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file = args.file
    if file is None or file == ".":
        file = ""
    print(FileUtils.file_extension(file))

def filename():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file = args.file
    if file is None or file == ".":
        file = ""
    print(FileUtils.filename_without_ext(file))

def delete_file():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file = args.file
    if FileUtils.is_file(file):
        FileUtils.delete_file(file)
        if FileUtils.is_file(file):
            LoggerUtils.error_log(f"File '{file}' not deleted")
    else:
        LoggerUtils.error_log(f"Not found File: {file}")

def file_encoding():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    print(FileUtils.get_file_encoding(args.file))

def writefile(args_list: list[str]|None =None):
    encoding: str = "utf-8"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", metavar="FILEPATH", type=str, required=True)
    parser.add_argument("-c", "--content", metavar="CONTENT", nargs="?", default="")
    parser.add_argument("-m", "--mode", choices=["w", "a"], default="w", help="Write mode (DEFAULT: w)")
    parser.add_argument("-e", "--enconding", metavar="ENCONDING", type=str, default=encoding, help="Type of enconding (DEFAULT: utf-8)")
    parser.add_argument("-d", "--force-dir", action="store_true", dest="force_dir", help="Create directory if not exists")
    args = parser.parse_args(args_list)
    filepath: str = args.file
    content = args.content
    options: WriteFileOptions = WriteFileOptions(
        mode=args.mode if args.mode else 'w',
        encoding=args.enconding if args.enconding is not None else encoding,
        force_dir=args.force_dir
    )
    if len(filepath) == 0:
        LoggerUtils.warn_log(f"Invalid given file path: {filepath}")
    else:
        res = FileUtils.write_file(filepath, content, options)
        if res:
            LoggerUtils.ok_log("File writing process completed.")
        else:
            LoggerUtils.error_log("File writing process completed unsuccessfully.")

def delete_file_lines(args_list: list[str]|None = None):
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", metavar="FILEPATH", type=str, required=True)
    parser.add_argument("-m", "--match", metavar="MATCH", type=str, required=True, help="Match content in lines to delete")
    args = parser.parse_args(args_list)
    filepath: str = args.file
    match: str = args.match
    if FileUtils.is_file(filepath):
        FileUtils.delete_line_on_file(filepath, match)
    else:
        LoggerUtils.error_log(f"File '{filepath}' not found.")

def file_contain():
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", metavar="FILEPATH", type=str, required=True)
    parser.add_argument("-m", "--match", metavar="MATCH", type=str, required=True, help="Match content in file")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable search with case insensitive")
    args = parser.parse_args()
    file_name: str = args.file
    match: str = args.match
    case_insensitive: bool = args.case_insensitive
    print(LoggerUtils.get_bool_str_formated(FileUtils.file_contains(file_name, match, case_insensitive)))

def findfile():
    parser = argparse.ArgumentParser()
    parser.add_argument("name", help="File name to find")
    args = parser.parse_args()
    cwd = os.getcwd()
    if args.name is not None and len(args.name) > 0:
        if FileUtils.is_dir(cwd):
            file_list = FileUtils.get_list_files_on_folder(cwd)
            for file in file_list:
                basename = os.path.basename(file)
                if args.name in basename:
                    print(file)
    else:
        LoggerUtils.error_log(f"Invalid given name: {args.name}")

def open_markdown():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file = args.file
    if FileUtils.is_file(file):
        base_cmd = ""
        if SYSTEM_UTILS.is_windows:
            base_cmd = "" if shutil.which("mdview") else "Get-Content"
        elif SYSTEM_UTILS.is_linux:
            base_cmd = (
                "markdown_viewer" if shutil.which("markdown_viewer") else "cat"
            )
        if len(base_cmd) > 0:
            CONSOLE_UTILS.exec_real_time(CommandInfo(command=base_cmd, args=[file]))
    else:
        LoggerUtils.error_log(f"Invalid given file: {file}")

def open_image():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file = args.file
    if FileUtils.is_file(file):
        base_cmd = ""
        if SYSTEM_UTILS.is_windows:
            base_cmd = "Start-Process"
        elif SYSTEM_UTILS.is_linux:
            if shutil.which("eog"):
                base_cmd = "eog"
            else:
                LoggerUtils.error_log("Can not open given image. Command not found: eog")

        if len(base_cmd) > 0:
            CONSOLE_UTILS.exec_real_time(CommandInfo(command=base_cmd, args=[file]))
    else:
        LoggerUtils.error_log(f"Invalid given file: {file}")

def count_files():
    if SHELL_UTILS.is_powershell:
        cmd = "(Get-ChildItem -File -Recurse -Force | Select-Object -ExpandProperty FullName | Measure-Object).Count"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL, verbose=True, use_shell=True))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command="find . -type f | wc -l", verbose=True))

def lhf():
    if SHELL_UTILS.is_powershell:
        cmd = "Get-ChildItem -Path \"$pwd\" -File -Force | ForEach-Object { $_.FullName }"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, verbose=True, shell=EShell.POWERSHELL, use_shell=True))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command="find . -maxdepth 1 -type f", verbose=True))

def move_files_to_parent():
    if SHELL_UTILS.is_powershell:
        cmd = "Get-ChildItem -Path \"$pwd\" -Recurse -File -Force | Move-Item -Destination \"$pwd\" -Verbose"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, verbose=True, shell=EShell.POWERSHELL, use_shell=True))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command="find . -mindepth 2 -type f -print -exec mv {} . \\;", verbose=True))
