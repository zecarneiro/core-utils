import argparse

from coreutils.libs.const_lib import SYSTEM_UTILS, CONSOLE_UTILS, SHELL_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.generic_lib import read_config, write_config
from coreutils.libs.processors.alias_processor import AliasProcessor
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.processors.script_processor import ScriptProcessor
from coreutils.libs.processors.service_processor import ServiceProcessor
from coreutils.libs.processors.which_processor import WhichProcessor
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.libs.pythonutils.system_utils import SystemUtils
from coreutils.varfunctions import shell_profile_script


def system_dependencies_apps() -> list[str]:
    apps: list[str] = []
    if SYSTEM_UTILS.is_windows:
        apps.append("sudo")
        apps.append("gsudo")
    elif SYSTEM_UTILS.is_linux:
        apps.append("lsb_release")
    return apps

def process_post_install_for_system_function_file():
    script_processor(["install", "-n", "now", "-c", "date", "-s"])
    if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
        script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("pgrepc.ps1")])
        script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("pkillc.ps1")])
        if SYSTEM_UTILS.is_windows:
            script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("sudopwsh.ps1")])
            script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("gsudopwsh.ps1")])
            script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("sudopwsh-old.ps1")])
    if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.BASH]):
        if SHELL_UTILS.is_bash:
            script_processor(["install", "-n", "cls", "-c", "clear", "-s"])
            if SYSTEM_UTILS.is_linux:
                script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("pgrepc.sh")])
                script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("snap-uninstall.sh")])
                script_processor(["install", "-n", "update-menu-entries", "-c", "sudo update-desktop-database", "-s"])

def reload_shell_profile(only_msg: bool = False):
    current_shell = SHELL_UTILS.current_shell
    if current_shell not in [EShell.CMD, EShell.UNKNOWN]:
        shell_script = GenericUtils.capture_print_from_function(shell_profile_script)
        msg = "Reload current shell with"
        if SHELL_UTILS.is_powershell:
            msg = f"{msg}: . '{shell_script}'"
        else:
            msg = f"{msg}: source '{shell_script}'"
        if only_msg:
            print(msg)
        else:
            LoggerUtils.info_log(msg)

def prompt_style():
    available_value_list = [1, 2, 3, 4]
    value: int = 1
    msg_template = f"""The default value is {value}.
    The list of available styles are: {available_value_list}."""
    parser = argparse.ArgumentParser()
    parser.add_argument("-s", "--status", action="store_true", dest="status", help="Return current value of prompt style")  # store_true = DEFAULT False
    parser.add_argument("-v", "--value", metavar="STYLE_OPTION", type=int, default=value, choices=available_value_list, help=f"Change prompt style({msg_template})")
    args = parser.parse_args()
    status: bool = args.status
    value = args.value
    if not SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.BASH]):
        MessageProcessor.show_shell_msg([EShell.POWERSHELL, EShell.BASH], "prompt-style")
        return
    current_shell = SHELL_UTILS.current_shell
    shell_name: str = current_shell.value
    config_data = read_config()
    if status:
        print(f"{config_data.promptStyle.get(current_shell.value)}")
    else:
        config_data.promptStyle[shell_name] = value
        write_config(config_data)
        reload_shell_profile()

def is_admin():
    print(LoggerUtils.get_bool_str_formated(SYSTEM_UTILS.is_admin()))

def rebootc():
    user_input = input("Will be restart PC. Continue(y/N)? ")
    if user_input == 'Y' or user_input == 'y':
        cmd_info = CommandInfo(command=ConsoleUtils.which("shutdown"), verbose=True)
        if SYSTEM_UTILS.is_windows:
            cmd_info.args = ["/r", "/t", "0"]
            CONSOLE_UTILS.exec_real_time(cmd_info)
        elif SYSTEM_UTILS.is_linux:
            cmd_info.args = ["-r", "now"]
            CONSOLE_UTILS.exec_real_time(cmd_info)

def shutdownc():
    user_input = input("Will be shutdown PC. Continue(y/N)? ")
    if user_input == 'Y' or user_input == 'y':
        cmd_info = CommandInfo(command=ConsoleUtils.which("shutdown"), verbose=True)
        if SYSTEM_UTILS.is_windows:
            cmd_info.args = ["/s", "/t", "0"]
            CONSOLE_UTILS.exec_real_time(cmd_info)
        elif SYSTEM_UTILS.is_linux:
            cmd_info.args = ["-h", "now"]
            CONSOLE_UTILS.exec_real_time(cmd_info)

def command_exists():
    parser = argparse.ArgumentParser()
    parser.add_argument("command", type=str)
    args = parser.parse_args()
    command = args.command
    if SYSTEM_UTILS.is_windows or SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
        cmd = f"if (Get-Command {command} -ErrorAction SilentlyContinue)"
        cmd = cmd + " { Write-Output true } else { Write-Output false }"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd))
    elif SYSTEM_UTILS.is_linux:
        cmd = f"if command -v {command} >/dev/null 2>&1; then echo true; else echo false; fi"
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd))
    else:
        print(CONST.UNKNOWN)

def whichc(args_list: list[str] | None = None):
    parser = argparse.ArgumentParser()
    parser.add_argument("command", type=str)
    args = parser.parse_args(args_list)
    which_processor_lib = WhichProcessor(command=args.command)
    print(which_processor_lib.find_command())

def evalc(args_list: list[str] | None = None):
    parser = argparse.ArgumentParser()
    parser.add_argument("command", type=str)
    args = parser.parse_args(args_list)
    command = args.command
    ConsoleUtils.exec_by_system(CommandInfo(command=command, verbose=True))

def script_processor(args_list: list[str] | None = None):
    script_lib = ScriptProcessor()
    parser = argparse.ArgumentParser()
    parser.description = "Set/Unset scripts as command."
    parser = argparse.ArgumentParser()
    sub_parser = parser.add_subparsers(dest="command", required=True)
    # INSTALL
    install_parser = sub_parser.add_parser("install")
    install_parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    install_parser.add_argument("-c", "--content", metavar="CONTENT", type=str, required=True, help="Content for script")
    install_parser.add_argument("-s", "--shell-tag", action="store_true", dest="shell_tag", help="Insert shell tag. Ex: #!/usr/bin/env bash")
    # INSTALL FROM FILE
    install_file_parser = sub_parser.add_parser("install-file")
    install_file_parser.add_argument("-f", "--file", metavar="SCRIPT_FILEPATH", required=True, help="Install from given script file.")
    # UNINSTALL
    uninstall_parser = sub_parser.add_parser("uninstall")
    uninstall_parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    # LIST
    list_parser = sub_parser.add_parser("list")
    list_parser.add_argument("-f", "--filter", metavar="FILTER_SEARCH", type=str, help="Filter to search")
    args = parser.parse_args(args_list)
    match args.command:
        case "install":
            script_lib.install(args.name, args.content, args.shell_tag)
        case "install-file":
            script_lib.install_from_file(args.file)
        case "uninstall":
            script_lib.uninstall(args.name)
        case "list":
            LoggerUtils.title_log("List of all script commands")
            for i, part in enumerate(script_lib.get_all(args.filter), start=1):
                print(f"{i}. {part}")
        case _:
            parser.print_help()

def alias_processor(args_list: list[str]|None = None):
    alias_lib: AliasProcessor = AliasProcessor()
    parser = argparse.ArgumentParser()
    parser.description = "Add/Remove alias."
    parser = argparse.ArgumentParser()
    sub_parser = parser.add_subparsers(dest="command", required=True)
    # ADD
    add_parser = sub_parser.add_parser("add")
    add_parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    add_parser.add_argument("-c", "--content", metavar="CONTENT", nargs="?", required=True, help="Content of alias. Only work with --add arg")
    # DELETE
    delete_parser = sub_parser.add_parser("delete")
    delete_parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    delete_parser.add_argument("-s", "--system", action="store_true", dest="system", help="If is to delete Windows SO system alias")
    # LIST
    list_parser = sub_parser.add_parser("list")
    list_parser.add_argument("-f", "--filter", metavar="FILTER_SEARCH", type=str, help="Filter to search")
    args = parser.parse_args(args_list)
    match args.command:
        case "add":
            alias_lib.add(args.name, args.content)
        case "delete":
            alias_lib.delete(args.name, args.system)
        case "list":
            LoggerUtils.title_log("List of all alias")
            for i, part in enumerate(alias_lib.get_all(args.filter), start=1):
                print(f"{i}. {part}")
        case _:
            parser.print_help()

def service_processor():
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    parser.add_argument("-c", "--command", metavar="COMMAND", type=str, help="Only work with --add arg")
    parser.add_argument("-d", "--description", metavar="DESCRIPTION", type=str)
    parser.add_argument("-i", "--install", action="store_true", dest="install", help="Process installation of service")
    parser.add_argument("-u", "--uninstall", action="store_true", dest="uninstall", help="Process uninstallation of service")
    args = parser.parse_args()
    service_processor_lib = ServiceProcessor(name=args.name, command=args.command, description=args.description)
    if service_processor_lib.is_valid(True if args.uninstall else False):
        return
    elif args.install:
        service_processor_lib.install()
    elif args.uninstall:
        service_processor_lib.uninstall()
    else:
        parser.print_help()

def envc(args_list: list[str]|None = None):
    parser = argparse.ArgumentParser()
    sub_parser = parser.add_subparsers(dest="command", required=True)
    # EXISTS
    exists_parser = sub_parser.add_parser("exists")
    exists_parser.add_argument("name", type=str, help="Name of env var")
    # EXIST VALUE
    exist_value_parser = sub_parser.add_parser("exist-value")
    exist_value_parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    exist_value_parser.add_argument("-v", "--value", metavar="VALUE", type=str, required=True)
    # PRINT VALUES
    print_values_parser = sub_parser.add_parser("print-values")
    print_values_parser.add_argument("-n", "--name", type=str, metavar="NAME", required=True, help="Name of env var")
    # LIST
    sub_parser.add_parser("list")
    # LIST WITH VALUES
    sub_parser.add_parser("list-with-values")
    args = parser.parse_args(args_list)
    match args.command:
        case "exists":
            print(LoggerUtils.get_bool_str_formated(SystemUtils.env_var_exists(args.name)))
        case "exist-value":
            print(LoggerUtils.get_bool_str_formated(SystemUtils.env_var_has_value(args.name, args.value)))
        case "print-values":
            values = SystemUtils.env_var_values_as_list(args.name)
            print(f"{args.name} has {len(values)} entries:")
            for i, part in enumerate(values, start=1):
                print(f"{i}. {part}")
        case "list":
            count = 1
            for name, _ in SystemUtils.env_var_list().items():
                print(f"{count}. {name}")
                count += 1
        case "list-with-values":
            for name, _ in SystemUtils.env_var_list().items():
                LoggerUtils.separator_log(20)
                envc(["print-values", name])
                print(CONST.EOF)
        case _:
            parser.print_help()

def printenvc():
    envc(["list-with-values"])

def ver():
    cmd = "lsb_release -a"
    if SYSTEM_UTILS.is_windows:
        cmd = "systeminfo | findstr /B /C:\"OS Name\" /B /C:\"OS Version\""
    ConsoleUtils.exec_by_system(CommandInfo(command=cmd))

def restart_explorer():
    command_info = CommandInfo(command="nautilus -q")
    if SYSTEM_UTILS.is_windows:
        command_info.command = "Stop-Process -Name explorer -Force; Start-Process explorer.exe"
        command_info.shell = EShell.POWERSHELL
        CONSOLE_UTILS.exec_real_time(command_info)
    else:
        LoggerUtils.info_log("Only works with nautilus")
        ConsoleUtils.exec_by_system(command_info)

def trash():
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)
    args = parser.parse_args()
    file = args.file
    if GenericUtils.str_is_empty(file) or file == ".":
        file = ""
    if not FileUtils.file_exist(file):
        LoggerUtils.error_log(f"File not found: {file}")
        return
    if SYSTEM_UTILS.is_windows:
        script_file = FileUtils.resolve_path(f"{DirsLib.get_resource_shell_script_libs_dir()}/trash.ps1")
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=f". \"{script_file}\" \"{file}\"", shell=EShell.POWERSHELL))
    else:
        command = f"mv --force -t ~/.local/share/Trash/files \"{file}\""
        ConsoleUtils.exec_by_system(CommandInfo(command=command))
