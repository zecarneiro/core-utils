import argparse

from coreutils.libs.const_lib import SYSTEM_UTILS, CONSOLE_UTILS, SHELL_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.processors.script_updater_processor import ScriptUpdaterProcessor
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.platform_enum import EPlatform
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.systemfunctions import script_processor


def package_dependencies_apps() -> list[str]:
    apps: list[str] = ["nodejs"]
    if SYSTEM_UTILS.is_windows:
        apps.append("winget")
        apps.append("scoop")
        apps.append("wsl")
    elif SYSTEM_UTILS.is_linux:
        apps.append("apt")
        apps.append("flatpak")
        apps.append("snap")
        apps.append("deb-get")
        apps.append("pacstall")
        apps.append("topgrade")
    return apps

def process_post_install_for_package_function_file():
    if SYSTEM_UTILS.is_linux and SHELL_UTILS.is_bash:
        script_processor(["install-file", "-f", DirsLib.get_resource_shell_script_apps_dir("snap-clean.sh")])

# ---------------------------------------------------------------------------- #
#                                      NPM                                     #
# ---------------------------------------------------------------------------- #
def npm_list():
    command_to_run: str = "npm list"
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--local", action="store_true", dest="local", help="List local package") # store_true = DEFAULT False
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH_NAME", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    is_local: bool = args.local
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_to_run = f"{command_to_run} --depth=0" if is_local else f"{command_to_run} -g --depth=0"
    command_info = CommandInfo(command=command_to_run, use_shell=True)
    if filter_app is None or len(filter_app) == 0:
        CONSOLE_UTILS.exec_real_time(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

def npm_clean(args_list: list[str]|None = None):
    LoggerUtils.title_log("CLEANUP NPM")
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--local", action="store_true", dest="local", help="Clean for local package")  # store_true = DEFAULT False
    args = parser.parse_args(args_list)
    is_local: bool = args.local
    cmd = f"npm {"-g" if not is_local else ""} cache clean --force"
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, verbose=True, use_shell=True))

# ---------------------------------------------------------------------------- #
#                                    WINGET                                    #
# ---------------------------------------------------------------------------- #
def winget_uninstall():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "winget-uninstall")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("app", type=str, help="App name")
    args = parser.parse_args()
    ConsoleUtils.exec_by_system(CommandInfo(command=f"winget uninstall --purge {args.app}", verbose=True))

def winget_list():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "winget-list")
        return
    command_to_run: str = "winget list"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

# ---------------------------------------------------------------------------- #
#                                     SCOOP                                    #
# ---------------------------------------------------------------------------- #
def scoop_uninstall():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "scoop-uninstall")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("app", type=str, help="App name")
    args = parser.parse_args()
    ConsoleUtils.exec_by_system(CommandInfo(command=f"scoop uninstall --purge {args.app}", verbose=True))

def scoop_list():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "scoop-list")
        return
    command_to_run: str = "scoop list"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

def scoop_clean():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "scoop-clean")
        return
    LoggerUtils.title_log("CLEANUP SCOOP")
    for cmd in ["scoop cleanup --all", "scoop cache rm *"]:
        ConsoleUtils.exec_by_system(CommandInfo(command=cmd, verbose=True))

# ---------------------------------------------------------------------------- #
#                                      WSL                                     #
# ---------------------------------------------------------------------------- #
def wsl_uninstall():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "wsl-uninstall")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("distro", type=str, help="Distro to uninstall")
    args = parser.parse_args()
    ConsoleUtils.exec_by_system(CommandInfo(command=f"wsl --unregister {args.distro}", verbose=True))

def wsl_list():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "wsl-list")
        return
    command_to_run: str = "wsl --list --verbose"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Distro to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

def wsl_shutdown():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "wsl-shutdown")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--force", action="store_true", dest="force")  # store_true = DEFAULT False
    args = parser.parse_args()
    cmd = "wsl --shutdown"
    if args.force:
        cmd = "sudo taskkill /F /IM wslservice.exe"
    ConsoleUtils.exec_by_system(CommandInfo(command=cmd))

def wsl_configc():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "wsl-configc")
        return
    config_file = FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}/.wslconfig")
    parser = argparse.ArgumentParser()
    parser.description = "This configurations only works on windows 11 or newer!!"
    parser.add_argument("-r", "--ram", metavar="MAX_RAM", type=int, help="Max RAM(GB) that WSL will use")
    parser.add_argument("-p", "--processor", metavar="MAX_PROCESSOR", type=int, help="Max Processor that WSL will use")
    args = parser.parse_args()
    ram: int = args.ram
    processor: int = args.processor
    data: str = ""
    if FileUtils.is_file(config_file):
        backup = f"{config_file}.bk"
        LoggerUtils.warn_log(f"Config file already exists: {config_file}")
        LoggerUtils.info_log(f"Backup original config file to: {backup}")
        FileUtils.copy_file(config_file, backup)
    if ram:
        data = f"memory={data}{CONST.EOF}{ram}GB"
    if processor:
        data = f"processors={processor}" if len(data) == 0 else f"{data}{CONST.EOF}processors={processor}"
    if len(data) > 0:
        FileUtils.write_file(config_file, f"[wsl2]{CONST.EOF}{data}")
        wsl_shutdown()

def wsl2win_path():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS], "wsl2win-path")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-p", "--path", metavar="PATH", type=str)
    parser.add_argument("-t", "--type", type=str, choices=["wsl2win", "win2wsl"], required=True)
    args = parser.parse_args()
    path: str = args.path
    cmd = "wsl -- wslpath {0} "'{1}'""
    if args.type == "wsl2win":
        cmd = cmd.format("-w", path)
    elif args.type == "win2wsl":
        cmd = cmd.format("-u", path)
    if not GenericUtils.str_is_empty(cmd):
        CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL))

# ---------------------------------------------------------------------------- #
#                                      APT                                     #
# ---------------------------------------------------------------------------- #
def apt_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "apt-uninstall")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("app", type=str, help="App name")
    args = parser.parse_args()
    ConsoleUtils.exec_by_system(CommandInfo(command=f"sudo apt purge --autoremove '{args.app}' -y", verbose=True))

def apt_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "apt-list")
        return
    command_to_run: str = "apt-mark showmanual"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

def apt_clean():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "apt-clean")
        return
    LoggerUtils.title_log("CLEANUP APT")
    for cmd in ["sudo apt clean -y", "sudo apt autoremove -y", "sudo apt autopurge -y", "sudo apt autoclean -y"]:
        ConsoleUtils.exec_by_system(CommandInfo(command=cmd, verbose=True))

# ---------------------------------------------------------------------------- #
#                                    FLATPAK                                   #
# ---------------------------------------------------------------------------- #
def flatpak_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "flatpak-uninstall")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("app", type=str, help="App name")
    args = parser.parse_args()
    ConsoleUtils.exec_by_system(CommandInfo(command=f"flatpak uninstall --delete-data '{args.app}' -y", verbose=True))

def flatpak_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "flatpak-list")
        return
    command_to_run: str = "flatpak list --columns=application"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

def flatpak_clean():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "flatpak-clean")
        return
    LoggerUtils.title_log("CLEANUP FLATPAK")
    for cmd in ["flatpak uninstall --unused -y", "sudo rm -rfv /var/tmp/flatpak-cache*"]:
        ConsoleUtils.exec_by_system(CommandInfo(command=cmd, verbose=True))

# ---------------------------------------------------------------------------- #
#                                     SNAP                                     #
# ---------------------------------------------------------------------------- #
def snap_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "snap-uninstall")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("app", type=str, help="App name")
    args = parser.parse_args()
    script_file = FileUtils.resolve_path(f"{DirsLib.get_resource_shell_script_libs_dir()}/snap-uninstall.sh")
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"'{script_file}' {args.app}", shell=SHELL_UTILS.current_shell, verbose=True, use_shell=True))

def snap_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "snap-list")
        return
    command_to_run: str = "snap list | awk 'NR >=2{print $1}'"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

# ---------------------------------------------------------------------------- #
#                                    DEB-GET                                   #
# ---------------------------------------------------------------------------- #
def deb_get_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "deb-get-uninstall")
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("apps", type=str, help="App name")
    args = parser.parse_args()
    ConsoleUtils.exec_by_system(CommandInfo(command=f"sudo deb-get purge {args.app}", verbose=True))

def deb_get_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "deb-get-list")
        return
    command_to_run: str = "deb-get list | grep installed | grep -v deb-get | awk '{print $1}'"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

def deb_get_clean():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "deb-get-clean")
        return
    LoggerUtils.title_log("CLEANUP DEB-GET")
    ConsoleUtils.exec_by_system(CommandInfo(command="sudo deb-get clean", verbose=True))

# ---------------------------------------------------------------------------- #
#                                   PACSTALL                                   #
# ---------------------------------------------------------------------------- #
def pacstall_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX], "pacstall-list")
        return
    command_to_run: str = "pacstall -L | awk '{print $1}'"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str|None = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
    if filter_app is None or len(filter_app) == 0:
        command_info.verbose = True
        ConsoleUtils.exec_by_system(command_info)
    else:
        res = CONSOLE_UTILS.exec(command_info)
        if res.has_error():
            res.log_error()
        else:
            matches: list[str]
            if is_case_insensitive:
                matches = GenericUtils.grep_e(res.stdout, filter_app, False)
            else:
                matches = GenericUtils.grep_e(res.stdout, filter_app)
            for match in matches:
                print(match)

# ---------------------------------------------------------------------------- #
#                                    OTHERS                                    #
# ---------------------------------------------------------------------------- #
def script_updater_processor(args_list: list[str]|None = None):
    script_updater_lib = ScriptUpdaterProcessor()
    parser = argparse.ArgumentParser()
    sub_parser = parser.add_subparsers(dest="command", required=True)
    # INSTALL
    install_parser = sub_parser.add_parser("install")
    install_parser.add_argument("-f", "--file", metavar="SCRIPT_FILEPATH", type=str, help="Install shell script by given file script")
    # UNINSTALL
    uninstall_parser = sub_parser.add_parser("uninstall")
    uninstall_parser.add_argument("-n", "--name", metavar="SCRIPT_NAME", type=str, help="Uninstall shell script by given file script")
    # RUN
    run_parser = sub_parser.add_parser("run")
    run_parser.add_argument("-n", "--name", metavar="SCRIPT_NAME", type=str, help="Process specific script")
    run_parser.add_argument("-a", "--all", action="store_true", dest="run_all", help="Process all script. Priority over name")
    # LIST
    list_parser = sub_parser.add_parser("list")
    list_parser.add_argument("-f", "--filter", metavar="FILTER_SEARCH", type=str, help="Filter to search")
    args = parser.parse_args(args_list)
    match args.command:
        case "install":
            script_updater_lib.install(args.file)
        case "uninstall":
            script_updater_lib.uninstall(args.name)
        case "list":
            LoggerUtils.title_log("List of all script updaters")
            for i, part in enumerate(script_updater_lib.get_all(args.filter), start=1):
                print(f"{i}. {part}")
        case "run":
            script_updater_lib.run(args.name, args.run_all)
        case _:
            parser.print_help()

def system_upgrade():
    cmd = "topgrade --cleanup --allow-root --skip-notify --yes --disable helm uv deb_get"
    if not SYSTEM_UTILS.is_windows:
        cmd = f"{cmd} powershell"
    ConsoleUtils.exec_by_system(CommandInfo(command=cmd, verbose=True))
    LoggerUtils.title_log("deb-get")
    for cmd_deb_get in ["sudo deb-get update", "sudo deb-get upgrade"]:
        ConsoleUtils.exec_by_system(CommandInfo(command=cmd_deb_get, verbose=True))
    script_updater_processor(["run", "--all"])

def system_cleanup():
    npm_clean()
    if SYSTEM_UTILS.is_windows:
        scoop_clean()
    elif SYSTEM_UTILS.is_linux:
        apt_clean()
        flatpak_clean()
        deb_get_clean()
        CONSOLE_UTILS.exec_real_time(CommandInfo(command="snap-clean"))