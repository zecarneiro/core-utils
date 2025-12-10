import argparse
import os

from coreutils.libs.const_lib import CONSOLE_UTILS, SHELL_UTILS, SYSTEM_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.processors.extract_processor import ExtractProcessor
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.entities.command_response import CommandResponse
from coreutils.libs.pythonutils.enums.platform_enum import EPlatform
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.systemfunctions import script_processor, alias_processor, evalc, whichc


def tools_dependencies_apps() -> list[str]:
    apps: list[str] = []
    if SYSTEM_UTILS.is_linux:
        apps.append("wget")
        apps.append("dos2unix")
    elif SYSTEM_UTILS.is_windows:
        apps.append("dos2unix")
    return apps

def process_post_install_for_tools_function_file():
    if SYSTEM_UTILS.is_windows and SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
        script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("win2wslpath.ps1")])
        script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("wsl2winpath.ps1")])
        script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("set-full-access.ps1")])
    if SYSTEM_UTILS.is_linux:
        alias_processor(["add", "-n", "sha1", "-c", "openssl sha1"])
        alias_processor(["add", "-n", "md5", "-c", "openssl md5"])
        alias_processor(["add", "-n", "sha256", "-c", "openssl sha256"])
        if not SHELL_UTILS.is_powershell:
            script_processor(["install", "-n", "pause", "-c", "echo -n \"Press [ENTER] to continue...: \"; read var_name", "-s"])
            script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("delete-broken-symlinks.sh")])

def text_is_empty():
    parser = argparse.ArgumentParser()
    parser.add_argument("text", type=str)
    args = parser.parse_args()
    text = args.text
    print(LoggerUtils.get_bool_str_formated(GenericUtils.str_is_empty(text)))

def count_args():
    print(ConsoleUtils.count_args())

def cutc():
    parser = argparse.ArgumentParser()
    parser.add_argument("text", type=str)
    parser.add_argument("-d", "--delimiter", type=str, required=True)
    parser.add_argument("-dr", "--direction", choices=["L", "R"], type=str, required=True)
    args = parser.parse_args()
    data: str = args.text
    delimiter: str = args.delimiter
    data_result = data.split(delimiter, 1)
    match args.direction:
        case "L":
            print(data_result[0])
        case "R":
            print(data_result[1])
        case _:
            print("")

def extract():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    parser.add_argument("-d", "--destination", type=str)
    args = parser.parse_args()
    extract_processor = ExtractProcessor(args.file, args.destination)
    extract_processor.process()

def has_internet():
    lib_name = "has-internet"
    res: CommandResponse
    if SHELL_UTILS.is_powershell:
        lib_path = DirsLib.get_resource_shell_script_libs_dir(f"{lib_name}.ps1")
        res = CONSOLE_UTILS.exec(CommandInfo(command=f". \"{lib_path}\"", shell=EShell.POWERSHELL))
    else:
        lib_path = DirsLib.get_resource_shell_script_libs_dir(f"{lib_name}.sh")
        res = CONSOLE_UTILS.exec(CommandInfo(command=f". \"{lib_path}\"", shell=EShell.BASH))
    if res.has_error():
        res.log_error()
    else:
        print(res.stdout)

def download():
    parser = argparse.ArgumentParser()
    parser.add_argument("-u", "--url", metavar="URL", type=str, required=True)
    parser.add_argument("-o", "--output-file", dest="output_file", metavar="OUTPUT_FILE", type=str, required=True)
    args = parser.parse_args()
    url: str = args.url
    file: str = args.output_file
    has_internet_conn  = GenericUtils.capture_print_from_function(has_internet)
    if len(url) == 0:
        LoggerUtils.error_log("Invalid given url")
        return
    if len(file) == 0:
        LoggerUtils.error_log("Invalid given output file")
        return
    if has_internet_conn == "true":
        LoggerUtils.info_log(f"Downloading from URL: {url}")
        if SYSTEM_UTILS.is_windows:
            lib_path = DirsLib.get_resource_shell_script_libs_dir("download.ps1")
            CONSOLE_UTILS.exec_real_time(CommandInfo(command=f". \"{lib_path}\" \"{url}\" -o \"{file}\"", shell=EShell.POWERSHELL))
        else:
            ConsoleUtils.exec_by_system(CommandInfo(command=f"wget -O \"{file}\" \"{url}\" -q --show-progress"))
    else:
        LoggerUtils.error_log("No Internet connection available")

def lhiden():
    cmd = 'ls -d .* --color=auto'
    if SYSTEM_UTILS.is_windows:
        cmd = "cmd /c dir \"$pwd\" /adh"
    ConsoleUtils.exec_by_system(CommandInfo(command=cmd))

def resolve_pathc():
    parser = argparse.ArgumentParser()
    parser.add_argument("path", type=str)
    args = parser.parse_args()
    print(FileUtils.resolve_path(args.path))

def dos2unix_recursive():
    parser = argparse.ArgumentParser()
    parser.add_argument("-e", "--extension", type=str, required=True)
    parser.add_argument("-t", "--type", type=str, choices=["dos2unix", "unix2dos"], required=True)
    args = parser.parse_args()
    extension: str = args.extension
    if not GenericUtils.str_is_empty(extension):
        if SYSTEM_UTILS.is_windows:
            cmd = f"Get-ChildItem -Recurse -File -Filter \"*.{extension}\"  | ForEach-Object "
            if args.type == "dos2unix":
                cmd = cmd + "{ dos2unix $_.FullName }"
            elif args.type == "unix2dos":
                cmd = cmd + "{ unix2dos $_.FullName }"
            CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL, verbose=True))
        else:
            cmd = f"find . -type f -name \"*.{extension}\" -print0 | xargs -0 "
            if args.type == "dos2unix":
                cmd = cmd + "dos2unix"
            elif args.type == "unix2dos":
                cmd = cmd + "unix2dos"
            ConsoleUtils.exec_by_system(CommandInfo(command=cmd, verbose=True))

def run_line_as_command():
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", type=str, required=True)
    parser.add_argument("-s", "--section", type=str, help="Section to run(ALL key included)")
    args = parser.parse_args()
    file: str = args.file
    section: str = args.section
    lib_name = "run-line-as-command"
    if SHELL_UTILS.is_powershell:
        cmd = f". {DirsLib.get_resource_shell_script_libs_dir(f"{lib_name}.ps1")} -file \"{file}\" -section \"{section}\""
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL))
    else:
        cmd = f". {DirsLib.get_resource_shell_script_libs_dir(f"{lib_name}.sh")} \"{file}\" \"{section}\""
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL))

def chmod_777():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file: str = args.file
    if SYSTEM_UTILS.is_windows:
        cmd = f"Unblock-File -Path \"{file}\""
        if FileUtils.is_dir(file):
            cmd = f"Get-ChildItem -Path \"{file}\" -Recurse | Unblock-File | Out-Null"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL))
    else:
        ConsoleUtils.exec_by_system(CommandInfo(command=f"chmod -R 777 \"{file}\""))

def unix2win_path():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "unix2win-path")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-p", "--path", metavar="PATH", type=str)
    parser.add_argument("-t", "--type", type=str, choices=["unix2win", "win2unix"], required=True)
    args = parser.parse_args()
    path: str = args.path
    cmd = ""
    if args.type == "unix2win":
        cmd = f"cygpath -w \"{path}\""
    elif args.type == "win2unix":
        cmd = f"cygpath -u \"{path}\""
    if not GenericUtils.str_is_empty(cmd):
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL))

def pyinit():
    file_path = FileUtils.resolve_path(f"{os.getcwd()}/__init__.py")
    FileUtils.touch(file_path)

def run_bin_or_cmd():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "run-bin")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("file_or_cmd", type=str)
    args = parser.parse_args()
    file_or_cmd: str = args.file_or_cmd
    command_info: CommandInfo|None = None
    if FileUtils.is_file(file_or_cmd) or :

        if not GenericUtils.str_is_empty(ConsoleUtils.which(file_or_cmd)):
            command_info = CommandInfo(command=cmd_start_process.format(file_or_cmd), shell=EShell.POWERSHELL)
        else:
            file_ext = FileUtils.file_extension(file_or_cmd)
            match file_ext:
                case ".exe" | ".msi":
                    command_info = CommandInfo(command=cmd_start_process.format(file_or_cmd), shell=EShell.POWERSHELL)
                case ".msixbundle" | ".msix":
                    command_info = CommandInfo(command=f"Add-AppxPackage -Path \"{file_or_cmd}\"", shell=EShell.POWERSHELL)
                case _:
                    LoggerUtils.error_log(f"Can not run this file: {file_or_cmd}")
                    command_info = None
    if command_info is not None:
        CONSOLE_UTILS.exec_real_time(command_info)
        LoggerUtils.info_log()

