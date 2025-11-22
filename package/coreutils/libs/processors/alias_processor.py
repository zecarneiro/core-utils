from dataclasses import dataclass

from coreutils.libs.const_lib import SYSTEM_UTILS, SHELL_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.entities.write_file_options import WriteFileOptions
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


@dataclass
class AliasProcessor:
    name: str
    content: str
    is_system: bool
    alias_file: str|None = None
    is_valid_shell: bool = True
    valid_shell: list[EShell]|None = None
    alias_dir: str = DirsLib.get_coreutils_shell_alias()
    is_alias_dir_created: bool = False

    def __post_init__(self):
        self.valid_shell = [EShell.POWERSHELL, EShell.BASH, EShell.KSH, EShell.ZSH]
        self.is_alias_dir_created = FileUtils.create_dir(self.alias_dir)
        if not SHELL_UTILS.is_shell(self.valid_shell):
            self.is_valid_shell = False
            MessageProcessor.show_shell_msg(self.valid_shell)
        if self.is_alias_dir_created:
            if SYSTEM_UTILS.is_powershell:
                self.alias_file = FileUtils.resolve_path(f"{self.alias_dir}/powershell.ps1")
            elif SYSTEM_UTILS.is_bash:
                self.alias_file = FileUtils.resolve_path(f"{self.alias_dir}/bash")
            elif SYSTEM_UTILS.is_ksh:
                self.alias_file = FileUtils.resolve_path(f"{self.alias_dir}/ksh")
            elif SYSTEM_UTILS.is_zsh:
                self.alias_file = FileUtils.resolve_path(f"{self.alias_dir}/zsh")
            else:
                LoggerUtils.error_log(f"{CONST.UNKNOWN} SHELL. Can not continue")
            if not FileUtils.is_file(self.alias_file) and not FileUtils.touch(self.alias_file):
                LoggerUtils.error_log(f"Failed on create alias file: {self.alias_file}")
                self.alias_file = None
        else:
             LoggerUtils.error_log(f"Failed on create dir: {self.alias_dir}")

    def __get_content_data_with_name(self) -> str|None:
        if SYSTEM_UTILS.is_powershell:
            return f"Set-Alias -Name \"{self.name}\" -Value"
        elif SYSTEM_UTILS.is_bash or SYSTEM_UTILS.is_ksh or SYSTEM_UTILS.is_zsh:
            return f"alias {self.name}="
        return None

    def is_valid(self, skip_content: bool = False) -> bool:
        if self.alias_file is None or not self.is_alias_dir_created or not self.is_valid_shell:
            return False
        elif GenericUtils.str_is_empty(self.name):
            LoggerUtils.error_log(f"Invalid given name: {self.name}")
            return False
        elif not GenericUtils.str_is_empty(self.name):
            if not skip_content and GenericUtils.str_is_empty(self.content):
                LoggerUtils.error_log(f"Invalid given content: {self.content}")
                return False
        return True

    def get_all(self) -> list[str]:
        alias_file_data = FileUtils.read_file(self.alias_file)
        alias_list: list[str] = []
        if alias_file_data is not None and len(alias_file_data) > 0:
            for alias_line in alias_file_data.splitlines():
                result = alias_line.split("=", 1)
                if result is not None and len(result) > 1:
                    alias_name = result[0].replace("alias ", "")
                    if alias_name is not None and len(alias_name) > 0:
                        alias_list.append(alias_name)
        return alias_list

    def __delete_process(self, verbose: bool):
        content_data: str | None = self.__get_content_data_with_name()
        if content_data is not None:
            FileUtils.delete_line_on_file(self.alias_file, content_data)
            if self.is_system:
                content_data = f"Remove-Item Alias:{self.name} -ErrorAction SilentlyContinue"
                FileUtils.write_file(self.alias_file, content_data, WriteFileOptions(mode="a"))
            if verbose:
                LoggerUtils.ok_log(f"Deleted alias with name: {self.name}.")
        else:
            if verbose:
                LoggerUtils.error_log(f"Failed on delete alias with name: {self.name}")

    def delete(self):
        self.__delete_process(True)

    def add(self):
        status: bool = False
        content_data = self.__get_content_data_with_name()
        if content_data is not None:
            self.__delete_process(False)
            if SYSTEM_UTILS.is_powershell:
                content_data = f"{content_data} \"{self.content}\""
                status = FileUtils.write_file(self.alias_file, content_data, WriteFileOptions(mode="a"))
            elif SYSTEM_UTILS.is_bash:
                content_data = f"{content_data}\"{self.content}\""
                status = FileUtils.write_file(self.alias_file, content_data, WriteFileOptions(mode="a"))
        if status:
            LoggerUtils.ok_log(f"Added alias with name: {self.name}")
        else:
            LoggerUtils.error_log(f"Failed on added alias with name: {self.name}")
