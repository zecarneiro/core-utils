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
    script_processor(["-i", "-n", "now", "-c", "date"])
    if SYSTEM_UTILS.is_windows:
        script_processor(["-i", "-n", "sudopwsh", "-c", f"sudo powershell -Command {CONST.POWERSHELL_ALL_ARGS_VAR_STR}"])
        script_processor(["-i", "-n", "gsudopwsh", "-c", f"gsudo powershell -Command {CONST.POWERSHELL_ALL_ARGS_VAR_STR}"])
        script_processor(["-i", "-n", "sudopwsh-old", "-c", f"Start-Process powershell.exe -verb runAs -Args \"{CONST.POWERSHELL_ALL_ARGS_VAR_STR}; pause\""])
        script_processor(["-i", "-n", "ver", "-c", "systeminfo | findstr /B /C:\"OS Name\" /B /C:\"OS Version\""])
        script_processor(["-f", f"{DirsLib.get_resource_shell_script_libs_dir()}/trash.ps1"])
        script_processor(["-i", "-n", "restart-explorer", "-c", "Stop-Process -Name explorer -Force; Start-Process explorer.exe"])
    else:
        script_processor(["-i", "-n", "ver", "-c", "lsb_release -a"])
        script_processor(["-i", "-n", "trash", "-c", f"mv --force -t ~/.local/share/Trash/files {SHELL_UTILS.get_all_args_var_name()}"])
        script_processor(["-i", "-n", "restart-explorer", "-c", "nautilus -q"])
        if not SHELL_UTILS.is_powershell:
            alias_processor(["-a", "-n", "cls", "-c", "clear"])
        script_processor(["-i", "-n", "update-menu-entries", "-c", "sudo update-desktop-database"])

    if SHELL_UTILS.is_powershell:
        script_processor(["-f", f"{DirsLib.get_resource_shell_script_libs_dir()}/pgrep.ps1"])
        script_processor(["-f", f"{DirsLib.get_resource_shell_script_libs_dir()}/pkill.ps1"])
    else:
        script_processor(["-i", "-n", "pgrep", "-c", f"{CONSOLE_UTILS.which("pgrep")} -l {CONST.BASH_ALL_ARGS_VAR_STR}"])

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
        MessageProcessor.show_shell_msg([EShell.POWERSHELL, EShell.BASH])
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

def reboot():
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
    parser.add_argument("-c", "--command", metavar="COMMAND", type=str, required=True, help="Command")
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

def whichc():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--command", metavar="COMMAND", type=str, required=True, help="Command")
    args = parser.parse_args()
    which_processor_lib = WhichProcessor(command=args.command)
    print(which_processor_lib.find_command())

def script_processor(args_list: list[str] | None = None):
    parser = argparse.ArgumentParser()
    parser.description = "Set scripts as command."
    parser.add_argument("-n", "--name", metavar="NAME", type=str)
    parser.add_argument("-c", "--content", metavar="CONTENT", type=str, help="Content for script")
    parser.add_argument("-i", "--install", action="store_true", dest="install", help="Process installation")
    parser.add_argument("-u", "--uninstall", action="store_true", dest="uninstall", help="Process uninstallation")
    parser.add_argument("-l", "--list", action="store_true", dest="list", help="List all of installed script")
    parser.add_argument("-f", "--file", metavar="SCRIPT_FILEPATH", help="Install from given script file. Install flag no need.")
    args = parser.parse_args(args_list)
    script_processor_lib: ScriptProcessor = ScriptProcessor(name=args.name, content=args.content, install_file=args.file)
    if args.list:
        LoggerUtils.title_log("List of all script commands")
        for file in script_processor_lib.get_all():
            print(f"- {file}")
    elif not script_processor_lib.is_valid(True if args.uninstall else False):
        return
    elif args.uninstall:
        script_processor_lib.uninstall()
    elif not GenericUtils.str_is_empty(script_processor_lib.install_file):
        script_processor_lib.install_from_file()
    elif args.install:
        script_processor_lib.install()
    else:
        LoggerUtils.warn_log("Please, set the flag --install|--uninstall|--list")
        parser.print_help()

def alias_processor(args_list: list[str]|None = None):
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--name", metavar="NAME", type=str)
    parser.add_argument("-c", "--content", metavar="CONTENT", nargs="?", help="Content of alias. Only work with --add arg")
    parser.add_argument("-a", "--add", action="store_true", dest="add", help="Process add alias")
    parser.add_argument("-d", "--delete", action="store_true", dest="delete", help="Process delete alias")
    parser.add_argument("-l", "--list", action="store_true", dest="list", help="List all of added alias")
    parser.add_argument("-s", "--system", action="store_true", dest="system", help="If is to delete Windows SO system alias")
    args = parser.parse_args(args_list)
    alias_processor_lib: AliasProcessor = AliasProcessor(name=args.name, content=args.content, is_system=args.system)
    if args.list:
        LoggerUtils.title_log("List of all alias")
        for alias_name in alias_processor_lib.get_all():
            print(f"- {alias_name}")
    elif not alias_processor_lib.is_valid(True if args.delete else False):
        return
    elif args.add:
        alias_processor_lib.add()
    elif args.delete:
        alias_processor_lib.delete()
    else:
        parser.print_help()

def evalc():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--command", metavar="COMMAND", type=str, required=True, help="Command")
    parser.add_argument("-s", "--shell", metavar="COMMAND", type=str, required=True, help="Command")
    args = parser.parse_args()
    command = args.command
    LoggerUtils.prompt_log(command)
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=command, shell=EShell.UNKNOWN))

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

def env_var_exists():
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    args = parser.parse_args()
    print(LoggerUtils.get_bool_str_formated(SystemUtils.env_var_exists(args.name)))

def env_var_has_value():
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    parser.add_argument("-v", "--value", metavar="VALUE", type=str, required=True)
    args = parser.parse_args()
    print(LoggerUtils.get_bool_str_formated(SystemUtils.env_var_has_value(args.name, args.value)))

def env_var_print_values(args_list: list[str]|None = None):
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--name", metavar="NAME", type=str, required=True)
    args = parser.parse_args(args_list)
    values = SystemUtils.env_var_values_as_list(args.name)
    print(f"{args.name} has {len(values)} entries:")
    for i, part in enumerate(values, start=1):
        print(f"{i}. {part}")

def printenvc():
    for name, _ in SystemUtils.env_var_list().items():
        LoggerUtils.separator_log(20)
        env_var_print_values(["-n", name])
        print(CONST.EOF)

