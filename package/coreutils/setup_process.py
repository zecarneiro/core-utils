import argparse
from importlib.metadata import distributions

from coreutils.directoryfunctions import process_post_install_for_dir_function_file
from coreutils.filefunctions import process_post_install_for_file_function_file, file_dependencies_apps
from coreutils.gitfunctions import process_post_install_for_git_function_file, git_dependencies_apps
from coreutils.libs.const_lib import SYSTEM_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.generic_lib import get_all_shell_profiles_files
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.entities.write_file_options import WriteFileOptions
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.packagefunctions import process_post_install_for_package_function_file, package_dependencies_apps, \
    script_updater_processor
from coreutils.systemfunctions import alias_processor, process_post_install_for_system_function_file, script_processor, \
    system_dependencies_apps
from coreutils.varfunctions import shell_profile_script, config_dir, other_apps_dir, user_bin_dir, user_startup_dir, \
    user_temp_dir, temp_dir, shell_script_dir, shell_alias_dir, shell_name

__HEADER_LENGTH__ = 35


def __get_all_package_cmds() -> list[str]:
    all_cmds: list[str] = []
    package_name = "coreutils"  # replace with your package name
    for dist in distributions():
        if dist.metadata["Name"].lower() == package_name.lower():
            eps = dist.entry_points
            for ep in eps:
                if ep.group == "console_scripts":
                    all_cmds.append(ep.name)
    return all_cmds

def __show_info_for_all_shell():
    LoggerUtils.info_log("If you have multiple shell in your machine, run this command on all")

def __setup_dependencies_apps() -> list[str]:
    apps: list[str] = []
    if SYSTEM_UTILS.is_windows:
        apps.append("bash")
        apps.append("gearlever")
    elif SYSTEM_UTILS.is_linux:
        apps.append("powershell")
    return apps

def __set_shell_script_and_alias():
    process_post_install_for_dir_function_file()
    process_post_install_for_file_function_file()
    process_post_install_for_git_function_file()
    process_post_install_for_package_function_file()
    process_post_install_for_system_function_file()
    script_processor(["-i", "-n", "gearlever", "-c", f"flatpak run it.mijorus.gearlever {SYSTEM_UTILS.get_all_args_var_name()}"])
    if SYSTEM_UTILS.is_windows:
        script_processor(["-i", "-n", "bash", "-c", f"$home\\scoop\\shims\\bash.exe {SYSTEM_UTILS.get_all_args_var_name()}"])
    else:
        if SYSTEM_UTILS.is_powershell:
            script_processor(["-i", "-n", "bash", "-c", f"{ConsoleUtils.which("bash")} {SYSTEM_UTILS.get_all_args_var_name()}"])
    if SYSTEM_UTILS.is_bash:
        alias_processor(["-a", "-n", "powershell", "-c", "pwsh -nologo"])
        alias_processor(["-a", "-n", "pwsh", "-c", "pwsh -nologo"])


def __install_shell_script():
    LoggerUtils.header_log("Create all shell scripts", __HEADER_LENGTH__)
    scripts: list[str] = [
        GenericUtils.capture_print_from_function(shell_profile_script),
    ]
    for script in scripts:
        if not FileUtils.is_file(script):
            print(script)
    LoggerUtils.ok_log("Done.")

def __create_default_dirs():
    LoggerUtils.header_log("Create all default directories", __HEADER_LENGTH__)
    directories: list[str] = [
        GenericUtils.capture_print_from_function(config_dir),
        GenericUtils.capture_print_from_function(other_apps_dir),
        GenericUtils.capture_print_from_function(user_bin_dir),
        GenericUtils.capture_print_from_function(user_startup_dir),
        GenericUtils.capture_print_from_function(user_temp_dir),
        GenericUtils.capture_print_from_function(temp_dir),
        GenericUtils.capture_print_from_function(shell_script_dir),
        GenericUtils.capture_print_from_function(shell_alias_dir)
    ]
    for directory in directories:
        if not FileUtils.is_dir(directory):
            print(directory)
    LoggerUtils.ok_log("Done.")

def __process_init_shell_script(is_enable: bool):
    init_script_file = "init-shell"
    if is_enable:
        LoggerUtils.header_log(f"Enable {init_script_file}", __HEADER_LENGTH__)
    else:
        LoggerUtils.header_log(f"Disable {init_script_file}", __HEADER_LENGTH__)
    for key, shell_file in get_all_shell_profiles_files().items():
        if FileUtils.is_file(shell_file):
            extension = None
            if key == EShell.POWERSHELL:
                extension = "ps1"
            elif key == EShell.BASH:
                extension = "sh"
            if extension is not None:
                coreutils_init_shell = FileUtils.resolve_path(f"{DirsLib.get_resource_dir()}/shell-scripts/{init_script_file}.{extension}")
                if is_enable:
                    if not FileUtils.file_contains(shell_file, coreutils_init_shell):
                        LoggerUtils.info_log(f"Inserting into shell script '{shell_file}' this core utils script: '{coreutils_init_shell}'")
                        options: WriteFileOptions = WriteFileOptions(
                            mode="a",
                            encoding=FileUtils.get_file_encoding(shell_file),
                        )
                        FileUtils.write_file(shell_file, f". \"{coreutils_init_shell}\"", options)
                else:
                    if FileUtils.file_contains(shell_file, coreutils_init_shell):
                        LoggerUtils.info_log(f"Deleting of shell script '{shell_file}' this core utils script: '{coreutils_init_shell}'")
                        FileUtils.delete_line_on_file(shell_file, coreutils_init_shell)
    if is_enable:
        LoggerUtils.ok_log(f"Please run this command in your terminal: . {GenericUtils.capture_print_from_function(shell_profile_script)}")
    else:
        LoggerUtils.ok_log(f"Please, restart your terminal")


def post_install():
    LoggerUtils.title_log(f"Running post install setup logic for SHELL: {GenericUtils.capture_print_from_function(shell_name)}")
    __show_info_for_all_shell()
    __create_default_dirs()
    __set_shell_script_and_alias()
    __install_shell_script()
    __process_init_shell_script(True)

def pre_uninstall():
    LoggerUtils.title_log(f"Running pre uninstall setup logic for SHELL: {GenericUtils.capture_print_from_function(shell_name)}")
    __show_info_for_all_shell()
    __process_init_shell_script(False)

def list_of_all_commands():
    LoggerUtils.title_log("List of all commands")
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="FILTER_TO_SEARCH", type=str, help="Install shell script by given file script")
    args = parser.parse_args()
    has_filter = True if not GenericUtils.str_is_empty(args.filter) else False
    for cmd in __get_all_package_cmds():
        if (has_filter and args.filter in cmd) or not has_filter:
            print(f"- {cmd}")
    print("")
    script_processor(["-l"])
    print("")
    alias_processor(["-l"])
    print("")
    script_updater_processor(["-l"])

def list_of_all_apps_dependencies():
    apps = __setup_dependencies_apps()
    apps = apps + file_dependencies_apps()
    apps = apps + git_dependencies_apps()
    apps = apps + package_dependencies_apps()
    apps = apps + system_dependencies_apps()
    LoggerUtils.info_log("Run this command to install most of all dependencies apps: coreutils-postinstall")
    LoggerUtils.title_log("List of all dependencies APPS")
    count = 1
    for app in apps:
        print(f"{count} - {app}")
        count += 1
