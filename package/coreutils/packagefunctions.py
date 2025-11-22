import argparse
import shutil

from coreutils.libs.const_lib import SYSTEM_UTILS, CONSOLE_UTILS, SHELL_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.generic_lib import set_file_permission_to_run
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.platform_enum import EPlatform
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
    if SYSTEM_UTILS.is_windows:
        script_processor(["-i", "-n", "scoop-clean", "-c", f"title-log -m 'CLEANUP SCOOP'{CONST.EOF}scoop cleanup --all; scoop cache rm *"])
    else:
        script_processor(["-i", "-n", "apt-clean", "-c", f"title-log -m 'CLEANUP APT'{CONST.EOF}sudo apt clean -y; sudo apt autoremove -y; sudo apt autopurge -y; sudo apt autoclean -y"])
        script_processor(["-i", "-n", "flatpak-clean", "-c", f"title-log -m 'CLEANUP FLATPAK'{CONST.EOF}flatpak uninstall --unused -y; sudo rm -rfv /var/tmp/flatpak-cache*"])
        script_processor(["-f", f"{DirsLib.get_resource_shell_script_libs_dir()}/snap-clean.sh"])
        script_processor(["-i", "-n", "deb-get-clean", "-c", f"title-log -m 'CLEANUP DEB-GET'{CONST.EOF}sudo deb-get clean"])

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
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_to_run = f"{command_to_run} --depth=0" if is_local else f"{command_to_run} -g --depth=0"
    command_info = CommandInfo(command=command_to_run)
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

def npm_clean(args_list=None):
    LoggerUtils.title_log("CLEANUP NPM")
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--local", action="store_true", dest="local", help="Clean for local package")  # store_true = DEFAULT False
    args = parser.parse_args(args_list)
    is_local: bool = args.local
    cmd = f"npm {"-g" if not is_local else ""} cache clean --force"
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, verbose=True))

# ---------------------------------------------------------------------------- #
#                                    WINGET                                    #
# ---------------------------------------------------------------------------- #
def winget_uninstall():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS])
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--app", metavar="APP", type=str, required=True, help="App name")
    args = parser.parse_args()
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"winget uninstall --purge {args.app}", verbose=True))

def winget_list():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS])
        return
    command_to_run: str = "winget list"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

# ---------------------------------------------------------------------------- #
#                                     SCOOP                                    #
# ---------------------------------------------------------------------------- #
def scoop_uninstall():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS])
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--app", metavar="APP", type=str, required=True, help="App name")
    args = parser.parse_args()
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"scoop uninstall --purge {args.app}", verbose=True))

def scoop_list():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS])
        return
    command_to_run: str = "scoop list"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

# ---------------------------------------------------------------------------- #
#                                      WSL                                     #
# ---------------------------------------------------------------------------- #
def wsl_list():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS])
        return
    command_to_run: str = "wsl --list --verbose"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Distro to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

def wsl_uninstall():
    if not SYSTEM_UTILS.is_windows:
        MessageProcessor.show_platform_msg([EPlatform.WINDOWS])
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--distro", metavar="DISTRO", type=str, required=True, help="Distro to uninstall")
    args = parser.parse_args()
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"wsl --unregister {args.distro}", verbose=True))

# ---------------------------------------------------------------------------- #
#                                      APT                                     #
# ---------------------------------------------------------------------------- #
def apt_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    command_to_run: str = "apt-mark showmanual"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

def apt_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--app", metavar="APP", type=str, required=True, help="App name")
    args = parser.parse_args()
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"sudo apt purge --autoremove '{args.app}' -y"))


# ---------------------------------------------------------------------------- #
#                                    FLATPAK                                   #
# ---------------------------------------------------------------------------- #
def flatpak_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    command_to_run: str = "flatpak list --columns=application"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

def flatpak_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--app", metavar="APP", type=str, required=True, help="App name")
    args = parser.parse_args()
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"flatpak uninstall --delete-data '{args.app}' -y"))


# ---------------------------------------------------------------------------- #
#                                     SNAP                                     #
# ---------------------------------------------------------------------------- #
def snap_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    command_to_run: str = "snap list | awk 'NR >=2{print \$1}'"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

def snap_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--app", metavar="APP", type=str, required=True, help="App name")
    args = parser.parse_args()
    script_file = FileUtils.resolve_path(f"{DirsLib.get_resource_shell_script_libs_dir()}/snap-uninstall.sh")
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"'{script_file}' {args.app}"))


# ---------------------------------------------------------------------------- #
#                                    DEB-GET                                   #
# ---------------------------------------------------------------------------- #
def deb_get_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    command_to_run: str = "deb-get list | grep installed | grep -v deb-get | awk '{print $1}'"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

def deb_get_uninstall():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--app", metavar="APP", type=str, required=True, help="App name")
    args = parser.parse_args()
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f"sudo deb-get purge {args.app}"))


# ---------------------------------------------------------------------------- #
#                                   PACSTALL                                   #
# ---------------------------------------------------------------------------- #
def pacstall_list():
    if not SYSTEM_UTILS.is_linux:
        MessageProcessor.show_platform_msg([EPlatform.LINUX])
        return
    command_to_run: str = "pacstall -L | awk '{print $1}'"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="APP_SEARCH", type=str, help="Package to search")
    parser.add_argument("-i", "--case-insensitive", action="store_true", dest="case_insensitive", help="Enable filter with case insensitive")  # store_true = DEFAULT False
    args = parser.parse_args()
    filter_app: str = args.filter
    is_case_insensitive: bool = args.case_insensitive
    command_info = CommandInfo(command=command_to_run)
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

# ---------------------------------------------------------------------------- #
#                                    OTHERS                                    #
# ---------------------------------------------------------------------------- #
def script_updater_processor(args_list=None):
    parser = argparse.ArgumentParser()
    parser.add_argument("-i", "--install", metavar="SCRIPT_FILEPATH", type=str, help="Install shell script by given file script")
    parser.add_argument("-u", "--uninstall", metavar="SCRIPT_NAME", type=str, help="Uninstall shell script by given file script")
    parser.add_argument("-r", "--run", metavar="NAME_SCRIPT|ALL", type=str, help="Process specific script or all")
    parser.add_argument("-l", "--list", action="store_true", dest="list", help="List all of installed script")
    parser.add_argument("-s", "--search", metavar="SCRIPT_TO_SEARCH", type=str, help="Search for installed script by given input")
    args = parser.parse_args(args_list)
    install_file: str = args.install
    uninstall_file: str = args.uninstall
    runner: str = args.run
    list_all: bool = args.list
    search_filter: str = args.search
    shell_path = FileUtils.resolve_path(f"{DirsLib.get_coreutils_local_dir()}/shell-scripts-installers-updaters/{SHELL_UTILS.current_shell.value}")
    file_list = FileUtils.get_list_files_on_folder(shell_path)
    if not GenericUtils.str_is_empty(install_file):
        if FileUtils.is_file(install_file) and FileUtils.create_dir(shell_path):
            dest = FileUtils.resolve_path(f"{shell_path}/{FileUtils.basename(install_file)}")
            can_install = True if not FileUtils.is_file(dest) else False
            if not can_install and ConsoleUtils.confirm("Script already exists. Continue", False):
                can_install = True
            if can_install:
                shutil.copy2(install_file, dest)
                set_file_permission_to_run(dest)
                LoggerUtils.ok_log(f"Installed shell script updater: {FileUtils.basename(install_file)}")
        else:
            LoggerUtils.error_log(f"Installation failed")
    elif not GenericUtils.str_is_empty(uninstall_file):
        for file in file_list:
            file_basename = FileUtils.basename(file)
            if uninstall_file == file_basename or uninstall_file == FileUtils.filename_without_ext(file_basename):
                FileUtils.delete_file(file)
                LoggerUtils.ok_log(f"Deleted  {file_basename}")
    elif not GenericUtils.str_is_empty(runner):
        is_all = runner == "ALL"
        if is_all:
            LoggerUtils.title_log(f"Process all installed shell script to install/update/uninstall packages")
        for file in file_list:
            file_basename = FileUtils.basename(file)
            can_run = True
            if not is_all:
                if runner == FileUtils.filename_without_ext(file_basename) or runner == file_basename:
                    can_run = True
                else:
                    can_run = False
            if can_run:
                LoggerUtils.header_log(f"Running {FileUtils.filename_without_ext(file_basename)}")
                CONSOLE_UTILS.exec_real_time(CommandInfo(command=f". \"{file}\""))
    elif not GenericUtils.str_is_empty(search_filter):
        for file in file_list:
            file_basename = FileUtils.basename(file)
            if search_filter in file_basename:
                print(f"- {file_basename}")
    elif list_all:
        LoggerUtils.title_log("List of all script updaters")
        for file in file_list:
            print(f"- {FileUtils.filename_without_ext(FileUtils.basename(file))}")

def system_upgrade():
    cmd = "topgrade --cleanup --allow-root --skip-notify --yes --disable helm uv deb_get"
    if not SYSTEM_UTILS.is_windows:
        cmd = f"{cmd} powershell"
    ConsoleUtils.exec_by_system(cmd, True)
    LoggerUtils.title_log("deb-get")
    CONSOLE_UTILS.exec_real_time(CommandInfo(command="sudo deb-get update && sudo deb-get upgrade", verbose=True))
    script_updater_processor(["-r", "ALL"])

def system_cleanup():
    npm_clean()
    if SYSTEM_UTILS.is_windows:
        CONSOLE_UTILS.exec_real_time(CommandInfo(command="scoop-clean"))
    elif SYSTEM_UTILS.is_linux:
        CONSOLE_UTILS.exec_real_time(CommandInfo(command="apt-clean"))
        CONSOLE_UTILS.exec_real_time(CommandInfo(command="flatpak-clean"))
        CONSOLE_UTILS.exec_real_time(CommandInfo(command="deb-get-clean"))
        CONSOLE_UTILS.exec_real_time(CommandInfo(command="snap-clean"))