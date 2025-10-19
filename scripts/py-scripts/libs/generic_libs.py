import sys

from vendor.pythonutils.console_utils import ConsoleUtils
from vendor.pythonutils.generic_utils import list_to_str

CONSOLE_UTILS = ConsoleUtils()

def get_args_list() -> list[str]:
    return sys.argv[1:]

def get_args_str() -> str:
    return list_to_str(sys.argv[1:])

def is_to_show_help(args: list[str], base_cmd: str, base_args: str) -> bool:
    if args is None or len(args) == 0:
        args = ["-h"]
    if "-h" in args or "--help" in args:
        print(f"usage: {base_cmd} [-h|--help]|[{base_args}]")
        return True
    return False

def install_alias_script():
    print("AAAA")