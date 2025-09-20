import argparse
import os
from importlib.metadata import distributions

from coreutils.filefunctions import file_dependencies_apps
from coreutils.gitfunctions import git_dependencies_apps
from coreutils.libs.const_lib import SYSTEM_UTILS, SHELL_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.generic_lib import get_all_shell_profiles_files
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.entities.write_file_options import WriteFileOptions
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.packagefunctions import process_post_install_for_package_function_file, package_dependencies_apps, \
    script_updater_processor
from coreutils.systemfunctions import alias_processor, process_post_install_for_system_function_file, \
    system_dependencies_apps, script_processor
from coreutils.toolsfunctions import tools_dependencies_apps, process_post_install_for_tools_function_file
from coreutils.varfunctions import shell_profile_script, config_dir, other_apps_dir, user_bin_dir, user_startup_dir, \
    user_temp_dir, temp_dir, shell_script_dir, shell_alias_dir, shell_name

__HEADER_LENGTH__ = 35


def __get_all_package_cmds(filter_name: str) -> list[str]:
    all_cmds: list[str] = []
    package_name = "coreutils"  # replace with your package name
    for dist in distributions():
        if dist.metadata["Name"].lower() == package_name.lower():
            eps = dist.entry_points
            for ep in eps:
                if ep.group == "console_scripts":
                    all_cmds.append(ep.name)
    if len(filter_name) > 0:
        all_cmds = [cmd for cmd in all_cmds if filter_name in cmd]
    return all_cmds

def __setup_dependencies_apps() -> list[str]:
    apps: list[str] = []
    if SYSTEM_UTILS.is_windows:
        apps.append("bash")
        apps.append("gearlever")
    elif SYSTEM_UTILS.is_linux:
        apps.append("powershell")
    return apps

def __set_shell_script_and_alias():
    process_post_install_for_package_function_file()
    process_post_install_for_system_function_file()
    process_post_install_for_tools_function_file()
    pwsh_bin = ConsoleUtils.which_by_cmds(["powershell", "pwsh"])
    alias_processor(["add", "-n", "powershell", "-c", f"{pwsh_bin} -nologo"])
    alias_processor(["add", "-n", "pwsh", "-c", f"{pwsh_bin} -nologo"])
    if SYSTEM_UTILS.is_windows:
        alias_processor(["add", "-n", "bash", "-c", "$home\\scoop\\shims\\bash.exe"])
    else:
        script_processor([
            "install",
            "-n", "gearlever",
            "-c", f"flatpak run it.mijorus.gearlever {SHELL_UTILS.get_all_args_var_name()}",
            "-s"
        ])

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
        GenericUtils.capture_print_from_function(shell_alias_dir),
    ]
    for directory in directories:
        if not FileUtils.is_dir(directory):
            print(directory)
    LoggerUtils.ok_log("Done.")

def __process_init_shell_script(is_enable: bool):
    shell_file = get_all_shell_profiles_files().get(SHELL_UTILS.current_shell)
    init_script_file = DirsLib.get_resource_shell_script_dir(f"init-{SHELL_UTILS.current_shell.value}-shell")
    how_to_reload_message = f". {GenericUtils.capture_print_from_function(shell_profile_script)}"
    match SHELL_UTILS.current_shell:
        case EShell.POWERSHELL:
            init_script_file = f"{init_script_file}.ps1"
        case EShell.BASH:
            init_script_file = f"{init_script_file}.sh"
        case EShell.KSH:
            init_script_file = f"{init_script_file}.ksh"
        case EShell.ZSH:
            init_script_file = f"{init_script_file}.zsh"
        case EShell.FISH:
            init_script_file = f"{init_script_file}.fish"
        case _:
            init_script_file = ""
    if not FileUtils.is_file(init_script_file) or (not FileUtils.is_file(shell_file) and not FileUtils.touch(shell_file)):
        return
    if is_enable:
        LoggerUtils.header_log(f"Enable {init_script_file}", __HEADER_LENGTH__)
        if not FileUtils.file_contains(shell_file, init_script_file):
            LoggerUtils.info_log("Inserting into:")
            print(f"- SHELL SCRIPT: {shell_file}")
            print(f"- THIS COREUTILS SCRIPT: {init_script_file}")
            options: WriteFileOptions = WriteFileOptions(mode="a")
            encoding = FileUtils.get_file_encoding(shell_file)
            if encoding is not None:
                options.encoding = encoding
            FileUtils.write_file(shell_file, f". \"{init_script_file}\"", options)
    else:
        LoggerUtils.header_log(f"Disable {init_script_file}", __HEADER_LENGTH__)
        if FileUtils.file_contains(shell_file, init_script_file):
            LoggerUtils.info_log("Deleting of:")
            print(f"- SHELL SCRIPT: {shell_file}")
            print(f"- THIS COREUTILS SCRIPT: {init_script_file}")
            FileUtils.delete_line_on_file(shell_file, init_script_file)
    LoggerUtils.ok_log(f"Please, restart your terminal: {how_to_reload_message}")

def __get_shell_installed() -> list[EShell]:
    shell_list: list[EShell] = []
    if SYSTEM_UTILS.is_windows or len(ConsoleUtils.which_by_cmds(["powershell", "pwsh"])) > 0:
        shell_list.append(EShell.POWERSHELL)
    if len(ConsoleUtils.which("bash")) > 0:
        shell_list.append(EShell.BASH)
    if len(ConsoleUtils.which("fish")) > 0:
        shell_list.append(EShell.FISH)
    if len(ConsoleUtils.which("ksh")) > 0:
        shell_list.append(EShell.KSH)
    if len(ConsoleUtils.which("zsh")) > 0:
        shell_list.append(EShell.ZSH)
    return shell_list

def post_install():
    for shell in __get_shell_installed():
        os.environ[CONST.PYTHON_UTILS_SHELL_NAME] = shell.value
        LoggerUtils.title_log(f"Running post install setup logic for SHELL: {GenericUtils.capture_print_from_function(shell_name)}")
        __create_default_dirs()
        __set_shell_script_and_alias()
        __install_shell_script()
        __process_init_shell_script(True)

def pre_uninstall():
    for shell in __get_shell_installed():
        os.environ[CONST.PYTHON_UTILS_SHELL_NAME] = shell.value
        LoggerUtils.title_log(f"Running pre uninstall setup logic for SHELL: {GenericUtils.capture_print_from_function(shell_name)}")
        __process_init_shell_script(False)

def list_of_all_commands():
    LoggerUtils.title_log("List of all commands")
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--filter", metavar="FILTER_TO_SEARCH", type=str)
    args = parser.parse_args()
    filter_name = args.filter if not GenericUtils.str_is_empty(args.filter) else ""
    for i, part in enumerate(__get_all_package_cmds(filter_name), start=1):
        print(f"{i}. {part}")
    print("")
    script_processor(["list", "--filter", filter_name])
    print("")
    alias_processor(["list", "--filter", filter_name])
    print("")
    script_updater_processor(["list", "--filter", filter_name])

def list_of_all_apps_dependencies():
    apps = __setup_dependencies_apps()
    apps = apps + file_dependencies_apps()
    apps = apps + git_dependencies_apps()
    apps = apps + package_dependencies_apps()
    apps = apps + system_dependencies_apps()
    apps = apps + tools_dependencies_apps()
    LoggerUtils.info_log("Run this command to install most of all dependencies apps: coreutils-postinstall")
    LoggerUtils.title_log("List of all dependencies APPS")
    count = 1
    for app in apps:
        print(f"{count} - {app}")
        count += 1
