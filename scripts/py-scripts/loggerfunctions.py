import argparse

from libs.generic_libs import get_args_list, is_to_show_help, get_args_str
from vendor.pythonutils.generic_utils import list_to_str
from vendor.pythonutils.logger_utils import info_log, error_log, debug_log, warn_log, ok_log, prompt_log, title_log, \
    header_log, separator_log, json_log

def log():
    args = get_args_list()
    if not is_to_show_help(args, "log", "MESSAGE"):
        print(get_args_str())

def errorlog():
    args = get_args_list()
    if not is_to_show_help(args, "errorlog", "MESSAGE"):
        error_log(get_args_str())

def infolog():
    args = get_args_list()
    if not is_to_show_help(args, "infolog", "MESSAGE"):
        info_log(get_args_str())

def debuglog():
    args = get_args_list()
    if not is_to_show_help(args, "debuglog", "MESSAGE"):
        debug_log(get_args_str())

def warnlog():
    args = get_args_list()
    if not is_to_show_help(args, "warnlog", "MESSAGE"):
        warn_log(get_args_str())

def oklog():
    args = get_args_list()
    if not is_to_show_help(args, "oklog", "MESSAGE"):
        ok_log(get_args_str())

def promptlog():
    args = get_args_list()
    if not is_to_show_help(args, "promptlog", "MESSAGE"):
        prompt_log(get_args_str())

def titlelog():
    fill_char: str = "#"
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", nargs="+", required=True, help="Title data")
    parser.add_argument("-c", "--char", type=str, default=fill_char, help="Character to fill the title. MUST BE 1 CHARACTER ONLY")
    args = parser.parse_args()
    message: str = list_to_str(args.message)
    if args.char is not None and len(args.char) == 1:
        fill_char = args.char
    title_log(message, fill_char)

def headerlog():
    fill_char: str = "-"
    length: int = 50
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", nargs="+", required=True, help="Header data")
    parser.add_argument("-l", "--length", type=int, default=length, help="Length of header")
    parser.add_argument("-c", "--char", type=str, default=fill_char, help="Character to fill the title. MUST BE 1 CHARACTER ONLY")
    args = parser.parse_args()
    message: str = list_to_str(args.message)
    if args.length is not None:
        length = args.length
    if args.char is not None and len(args.char) == 1:
        fill_char = args.char
    header_log(message, length, fill_char)

def separatorlog():
    fill_char: str = "-"
    length: int = 6
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--length", type=int, default=length, help="Length of header")
    parser.add_argument("-c", "--char", type=str, default=fill_char, help="Character to fill the title. MUST BE 1 CHARACTER ONLY")
    args = parser.parse_args()
    if args.length is not None:
        length = args.length
    if args.char is not None and len(args.char) == 1:
        fill_char = args.char
    separator_log(length, fill_char)

def jsonlog():
    args = get_args_list()
    if not is_to_show_help(args, "jsonlog", "JSON STRING"):
        json_log(get_args_str())
