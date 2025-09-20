import os
from dataclasses import dataclass

from coreutils.libs.const_lib import SYSTEM_UTILS, CONSOLE_UTILS
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


@dataclass
class ExtractProcessor:
    file: str|None
    destination: str|None

    def __post_init__(self):
        if self.destination is None or len(self.destination) == 0:
            self.destination = os.getcwd()
        if self.file is not None and len(self.file) == 0:
            self.file = None

    @property
    def __file_ext(self) -> str:
        return FileUtils.file_extension(self.file if self.file is not None else "", True)

    def __print_unknown_file(self):
        LoggerUtils.error_log(f"Don't know how to extract '{self.file}'...")

    def __get_cmd(self) -> str|None:
        match self.__file_ext:
            case ".zip":
                if SYSTEM_UTILS.is_windows:
                    return f"Expand-Archive -LiteralPath \"{self.file}\" -DestinationPath \"{self.destination}\""
                else:
                    return f"unzip \"{self.file}\" -d \"{self.destination}\""
            case ".tar.bz2" | ".tbz2":
                if SYSTEM_UTILS.is_linux:
                    return f"tar xvjf \"{self.file}\""
            case ".tar.gz" | ".tgz":
                if SYSTEM_UTILS.is_linux:
                    return f"tar xvzf \"{self.file}\""
            case ".bz2":
                if SYSTEM_UTILS.is_linux:
                    return f"bunzip2 \"{self.file}\""
            case ".rar":
                if SYSTEM_UTILS.is_linux:
                    return f"rar x \"{self.file}\""
            case ".gz":
                if SYSTEM_UTILS.is_linux:
                    return f"gunzip \"{self.file}\""
            case ".tar":
                if SYSTEM_UTILS.is_linux:
                    return f"tar xvf \"{self.file}\""
            case ".Z":
                if SYSTEM_UTILS.is_linux:
                    return f"uncompress \"{self.file}\""
            case ".7z":
                if SYSTEM_UTILS.is_linux:
                    return f"7z x \"{self.file}\""
            case _:
                self.__print_unknown_file()
        return None

    def __process_windows(self):
        cmd_str = self.__get_cmd()
        cmd_info: CommandInfo|None = None
        if cmd_str is not None:
            match self.__file_ext:
                case ".zip":
                    cmd_info = CommandInfo(command=cmd_str, use_shell=True)
                case _:
                    self.__print_unknown_file()
        if cmd_info is not None:
            CONSOLE_UTILS.exec_real_time(cmd_info)

    def __process_linux(self):
        cmd_str = self.__get_cmd()
        if cmd_str is not None:
            ConsoleUtils.exec_by_system(CommandInfo(command=cmd_str))

    def __is_valid(self) -> bool:
        if self.file is None or not FileUtils.is_file(self.file):
            LoggerUtils.error_log("Given file not found")
            return False
        return True

    def process(self):
        if self.__is_valid():
            if SYSTEM_UTILS.is_windows:
                self.__process_windows()
            else:
                self.__process_linux()
