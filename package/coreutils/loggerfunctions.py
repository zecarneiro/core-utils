import argparse

from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


def log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    args = parser.parse_args()
    print(GenericUtils.list_to_str(args.message))

def error_log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    args = parser.parse_args()
    LoggerUtils.error_log(GenericUtils.list_to_str(args.message))

def info_log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    args = parser.parse_args()
    LoggerUtils.info_log(GenericUtils.list_to_str(args.message))

def debug_log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    args = parser.parse_args()
    LoggerUtils.debug_log(GenericUtils.list_to_str(args.message))

def warn_log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    args = parser.parse_args()
    LoggerUtils.warn_log(GenericUtils.list_to_str(args.message))

def ok_log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    args = parser.parse_args()
    LoggerUtils.ok_log(GenericUtils.list_to_str(args.message))

def prompt_log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    args = parser.parse_args()
    LoggerUtils.prompt_log(GenericUtils.list_to_str(args.message))

def title_log():
    fill_char: str = "#"
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    parser.add_argument("-c", "--char", metavar="CHAR", type=str, default=fill_char, help="Character to fill the title. MUST BE 1 CHARACTER ONLY")
    args = parser.parse_args()
    message: str = GenericUtils.list_to_str(args.message)
    if args.char is not None and len(args.char) == 1:
        fill_char = args.char
    LoggerUtils.title_log(message, fill_char)

def header_log():
    fill_char: str = "-"
    length: int = 50
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", "--message", metavar="MESSAGE", nargs="+", required=True)
    parser.add_argument("-l", "--length", metavar="LENGTH", type=int, default=length, help="Length of header")
    parser.add_argument("-c", "--char", metavar="CHAR", type=str, default=fill_char, help="Character to fill the title. MUST BE 1 CHARACTER ONLY")
    args = parser.parse_args()
    message: str = GenericUtils.list_to_str(args.message)
    if args.length is not None:
        length = args.length
    if args.char is not None and len(args.char) == 1:
        fill_char = args.char
    LoggerUtils.header_log(message, length, fill_char)

def separator_log():
    fill_char: str = "-"
    length: int = 6
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--length", metavar="LENGTH", type=int, default=length, help="Length of header")
    parser.add_argument("-c", "--char", metavar="CHAR", type=str, default=fill_char, help="Character to fill the title. MUST BE 1 CHARACTER ONLY")
    args = parser.parse_args()
    if args.length is not None:
        length = args.length
    if args.char is not None and len(args.char) == 1:
        fill_char = args.char
    LoggerUtils.separator_log(length, fill_char)

def json_log():
    parser = argparse.ArgumentParser()
    parser.add_argument("-j", "--json", metavar="JSON_STRING", nargs="+", required=True)
    args = parser.parse_args()
    data: list[str] = args.json
    LoggerUtils.json_log(GenericUtils.list_to_str(data))
