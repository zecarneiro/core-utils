import argparse

from libs.generic_libs import CONSOLE_UTILS, get_args_list, is_to_show_help
from vendor.pythonutils.entities.command_info import CommandInfo
from vendor.pythonutils.enums.shell_enum import EShell
from vendor.pythonutils.system_utils import get_shell

def npmupgrade():
    CONSOLE_UTILS.exec_real_time(CommandInfo(command="npm", args=["outdated", "-g"]))
    CONSOLE_UTILS.exec_real_time(CommandInfo(command="npm", args=["update", "-g"]))

def npmlist():
    command_to_run: str = "npm list"
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--local", action="store_true", dest="local", help="List local package") # store_true = DEFAULT False
    parser.add_argument("-f", "--filter", type=str, help="Package to search")
    args = parser.parse_args()
    is_local: bool = args.local
    filter: str = args.filter
    command_to_run = f"{command_to_run} --depth=0" if is_local else f"{command_to_run} -g --depth=0"
    if filter is not None and len(filter) > 0:
        shell = get_shell()
        if shell == EShell.POWERSHELL or shell == EShell.CMD:
            command_to_run = f"{command_to_run} | Select-String \"{filter}\""
        else:
            command_to_run = f"{command_to_run} | grep \"{filter}\""
    CONSOLE_UTILS.exec_real_time(CommandInfo(command=command_to_run, verbose=True))

def npmclean():
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--local", action="store_true", dest="local", help="Clean for local package")  # store_true = DEFAULT False
    args = parser.parse_args()
    is_local: bool = args.local
    CONSOLE_UTILS.exec_real_time(CommandInfo(command="npm", args=["cache", "clean", "-g" if not is_local else "", "--force"]))

def installupdaterscript():
    args = get_args_list()
    if not is_to_show_help(args, "installupdaterscript", "FULL_UPDATER_SCRIPT_PATH"):
        print("AAA")


