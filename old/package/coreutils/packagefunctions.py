import argparse

from coreutils.libs.const_lib import CONSOLE_UTILS, SHELL_UTILS, SYSTEM_UTILS
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
        script_processor(
            [
                "install-file",
                "-f",
                DirsLib.get_resource_shell_script_apps_dir("snap-clean.sh"),
            ]
        )


# ---------------------------------------------------------------------------- #
#                                    OTHERS                                    #
# ---------------------------------------------------------------------------- #
def script_updater_processor(args_list: list[str] | None = None):
    script_updater_lib = ScriptUpdaterProcessor()
    parser = argparse.ArgumentParser()
    sub_parser = parser.add_subparsers(dest="command", required=True)
    # INSTALL
    install_parser = sub_parser.add_parser("install")
    install_parser.add_argument(
        "-f",
        "--file",
        metavar="SCRIPT_FILEPATH",
        type=str,
        help="Install shell script by given file script",
    )
    # UNINSTALL
    uninstall_parser = sub_parser.add_parser("uninstall")
    uninstall_parser.add_argument(
        "-n",
        "--name",
        metavar="SCRIPT_NAME",
        type=str,
        help="Uninstall shell script by given file script",
    )
    # RUN
    run_parser = sub_parser.add_parser("run")
    run_parser.add_argument(
        "-n", "--name", metavar="SCRIPT_NAME", type=str, help="Process specific script"
    )
    run_parser.add_argument(
        "-a",
        "--all",
        action="store_true",
        dest="run_all",
        help="Process all script. Priority over name",
    )
    # LIST
    list_parser = sub_parser.add_parser("list")
    list_parser.add_argument(
        "-f", "--filter", metavar="FILTER_SEARCH", type=str, help="Filter to search"
    )
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
