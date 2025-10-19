import argparse
import os
import shutil

from libs.generic_libs import get_args_list, is_to_show_help
from vendor.pythonutils.console_utils import ConsoleUtils
from vendor.pythonutils.entities.command_info import CommandInfo
from vendor.pythonutils.enums.shell_enum import EShell
from vendor.pythonutils.file_utils import file_exist, file_extension, filename_without_ext, delete_file, \
    get_file_encoding, file_contains, get_list_files_on_folder, get_list_dirs_on_folder
from vendor.pythonutils.logger_utils import get_bool_str_formated, info_log, error_log
from vendor.pythonutils.system_utils import is_windows, is_linux, get_shell

from vendor.pythonutils.generic_utils import list_to_str


def fileexists():
    args = get_args_list()
    if not is_to_show_help(args, "fileexists", "FILE"):
        print(get_bool_str_formated(file_exist(list_to_str(args))))

def fileextension():
    args = get_args_list()
    if not is_to_show_help(args, "fileextension", "FILE"):
        print(file_extension(list_to_str(args)))

def filename():
    args = get_args_list()
    if not is_to_show_help(args, "filename", "FILE"):
        print(filename_without_ext(list_to_str(args)))

def deletefile():
    args = get_args_list()
    if not is_to_show_help(args, "deletefile", "FILE"):
        file = list_to_str(args)
        if file_exist(file):
            delete_file(file)
            if not file_exist(file):
                info_log(f"Deleted file: {file}")
            else:
                error_log(f"File '{file}' not deleted")

def fileenconding():
    args = get_args_list()
    if not is_to_show_help(args, "deletefile", "FILE"):
        print(get_file_encoding(list_to_str(args)))

def writefile():
    encoding: str = "utf-8"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", type=str, required=True, help="/path/to/file")
    parser.add_argument("-c", "--contents", nargs="+", default=[""], help="Data to insert on file. NOTE: Example of string with break lines: \"aaa \\nbbb\"")
    parser.add_argument("-a", "--append", action="store_true", dest="append", help="Enable append to file") # store_true = DEFAULT False
    parser.add_argument("-e", "--enconding", type=str, default=encoding, help="Type of enconding")
    args = parser.parse_args()
    filename: str = args.file
    is_append: bool = args.append
    contents: list[str] = args.contents if args.contents is not None else []
    if args.enconding is not None:
        encoding = args.enconding
    write_type = "w"
    if len(filename) == 0:
        error_log(f"Invalid given file: {filename}")
    else:
        if file_exist(filename) and is_append:
            write_type = "a"
    with open(filename, write_type, encoding=encoding) as fstream:
        is_first = True
        for content in contents:
            if is_first:
                fstream.write(content)
                is_first = False
            else:
                fstream.write(f" {content}")

def delfilelines():
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", type=str, required=True, help="/path/to/file")
    parser.add_argument("-m", "--match", type=str, required=True, help="Match content in lines to delete")
    args = parser.parse_args()
    filename: str = args.file
    match: str = args.match
    if file_exist(filename):
        encoding = get_file_encoding(filename)
        with open(filename, "r", encoding=encoding) as fstream:
            lines = fstream.readlines()

        # Keep non-empty lines
        lines = [line for line in lines if match not in line]

        with open(filename, "w", encoding=encoding) as fstream:
            fstream.writelines(lines)

def countfiles():
    console_utils = ConsoleUtils()
    shell = get_shell()
    if shell == EShell.POWERSHELL or shell == EShell.CMD:
        console_utils.exec_real_time(CommandInfo(command="(Get-ChildItem -File -Recurse -Force | Select-Object -ExpandProperty FullName | Measure-Object).Count", use_shell=True))
    else:
        console_utils.exec_real_time(CommandInfo(command="find . -type f | wc -l"))

def movefilestoparent():
    console_utils = ConsoleUtils()
    shell = get_shell()
    if shell == EShell.POWERSHELL or shell == EShell.CMD:
        console_utils.exec_real_time(CommandInfo(command="Get-ChildItem -Path \"$pwd\" -Recurse -File -Force | Move-Item -Destination \"$pwd\" -Verbose", use_shell=True))
    else:
        console_utils.exec_real_time(CommandInfo(command="find . -mindepth 2 -type f -print -exec mv {} . \;"))

def lf():
    console_utils = ConsoleUtils()
    shell = get_shell()
    if shell == EShell.POWERSHELL or shell == EShell.CMD:
        console_utils.exec_real_time(CommandInfo(command="Get-ChildItem -Path \"$pwd\" -File -Force | ForEach-Object { $_.FullName }", use_shell=True))
    else:
        console_utils.exec_real_time(CommandInfo(command="find . -maxdepth 1 -type f"))

def filecontain():
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", type=str, required=True, help="/path/to/file")
    parser.add_argument("-m", "--match", type=str, required=True, help="Match content in file")
    parser.add_argument("--case-insensitive", action="store_true", dest="case_insensitive", help="Enable search with case insensitive") # store_true = DEFAULT False
    args = parser.parse_args()
    filename: str = args.file
    match: str = args.match
    case_insensitive: bool = args.case_insensitive
    print(get_bool_str_formated(file_contains(filename, match, case_insensitive)))

def findfile():
    args = get_args_list()
    if not is_to_show_help(args, "findfile", "FILENAME"):
        cwd = os.getcwd()
        if file_exist(cwd):
            file_list = get_list_files_on_folder(cwd)
            filename = list_to_str(args)
            for file in file_list:
                basename = os.path.basename(file)
                if filename in basename:
                    print(file)

def finddir():
    args = get_args_list()
    if not is_to_show_help(args, "finddir", "DIRNAME"):
        cwd = os.getcwd()
        if file_exist(cwd):
            dir_list = get_list_dirs_on_folder(cwd)
            dirname = list_to_str(args)
            for dir in dir_list:
                basedir = os.path.dirname(dir)
                if dirname in basedir:
                    print(dir)

# DEPENDS APP:
# - mdview
def openmarkdown():
    args = get_args_list()
    if not is_to_show_help(args, "openmarkdown", "FILE"):
        file = list_to_str(args)
        if file_exist(file):
            base_cmd = ""
            if is_windows():
                base_cmd = "mdview" if shutil.which("mdview") else "Get-Content"
            elif is_linux():
                base_cmd = "markdown_viewer" if shutil.which("markdown_viewer") else "cat"
            if len(base_cmd) > 0:
                console_utils = ConsoleUtils()
                console_utils.exec_real_time(CommandInfo(command=base_cmd, args=[file]))
        else:
            error_log(f"Invalid given file: {file}")

# DEPENDS APP:
# - On linux: eog
def openimage():
    args = get_args_list()
    if not is_to_show_help(args, "openimage", "FILE"):
        file = list_to_str(args)
        if file_exist(file):
            base_cmd = ""
            if is_windows():
                base_cmd = "Start-Process"
            elif is_linux():
                if shutil.which("eog"):
                    base_cmd = "eog"
                else:
                    error_log("Can not open given image. Command not found: eog")

            if len(base_cmd) > 0:
                console_utils = ConsoleUtils()
                console_utils.exec_real_time(CommandInfo(command=base_cmd, args=[file]))
        else:
            error_log(f"Invalid given file: {file}")


