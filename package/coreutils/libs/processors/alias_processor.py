from dataclasses import dataclass

from coreutils.libs.const_lib import SHELL_UTILS, SYSTEM_UTILS
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
    @property
    def __valid_shell(self) -> list[EShell]:
        return [EShell.POWERSHELL, EShell.CMD, EShell.BASH, EShell.KSH, EShell.ZSH, EShell.FISH]

    @property
    def __scripts_dir(self) -> str:
        if self.__is_valid_shell():
            directory = DirsLib.get_coreutils_shell_script()
            FileUtils.create_dir(directory)
            return directory
        return ""

    @property
    def __alias_dir(self) -> str:
        return DirsLib.get_coreutils_shell_alias()

    @property
    def __alias_file(self) -> str:
        alias_file = ""
        suffix = "-alias"
        if SHELL_UTILS.is_powershell:
            alias_file = FileUtils.resolve_path(f"{self.__alias_dir}/powershell{suffix}.ps1")
        elif SHELL_UTILS.is_bash:
            alias_file = FileUtils.resolve_path(f"{self.__alias_dir}/bash{suffix}")
        elif SHELL_UTILS.is_ksh:
            alias_file = FileUtils.resolve_path(f"{self.__alias_dir}/ksh{suffix}")
        elif SHELL_UTILS.is_zsh:
            alias_file = FileUtils.resolve_path(f"{self.__alias_dir}/zsh{suffix}")
        elif SHELL_UTILS.is_fish:
            alias_file = FileUtils.resolve_path(f"{self.__alias_dir}/fish{suffix}")
        if not FileUtils.is_file(alias_file):
            prefix = "#!/usr/bin/env"
            if SHELL_UTILS.is_bash:
                FileUtils.write_file(alias_file, f"{prefix} bash")
            elif SHELL_UTILS.is_zsh:
                FileUtils.write_file(alias_file, f"{prefix} zsh")
            elif SHELL_UTILS.is_ksh:
                FileUtils.write_file(alias_file, f"{prefix} ksh")
            elif SHELL_UTILS.is_fish:
                FileUtils.write_file(alias_file, f"{prefix} fish")
            else:
                FileUtils.touch(alias_file)
        return alias_file

    @property
    def __windows_cmd_content_template(self) -> str:
        return """
        @echo off
        powershell.exe -NoLogo -ExecutionPolicy Bypass -Command "{0} {1}" """

    def __is_valid_shell(self) -> bool:
        status = SHELL_UTILS.is_shell(self.__valid_shell)
        if not status:
            MessageProcessor.show_shell_msg(self.__valid_shell, self.__class__.__name__)
        return status

    def __is_alias_dir_created(self) -> bool:
        error_msg = "Failed on create dir"
        dir_exist = FileUtils.is_dir(self.__alias_dir)
        cmd_dir_exist = True
        if SYSTEM_UTILS.is_windows:
            cmd_dir = FileUtils.resolve_path(f"{self.__alias_dir}/{EShell.CMD.value}")
            cmd_dir_exist = FileUtils.is_dir(cmd_dir)
            if not cmd_dir_exist:
                LoggerUtils.error_log(f"{error_msg}: {cmd_dir}")
        if not dir_exist:
            LoggerUtils.error_log(f"{error_msg}: {self.__alias_dir}")
        return dir_exist and cmd_dir_exist

    def __is_alias_file_created(self) -> bool:
        status = True if len(self.__alias_file) > 0 and FileUtils.is_file(self.__alias_file) else False
        if not status:
            LoggerUtils.error_log(f"Failed on create alias file: {self.__alias_file}")
        return status

    def __is_valid(self) -> bool:
        return self.__is_valid_shell() and self.__is_alias_dir_created() and self.__is_alias_file_created()

    @staticmethod
    def __is_valid_name(name: str) -> bool:
        status = True if not GenericUtils.str_is_empty(name) else False
        if not status:
            LoggerUtils.error_log("Invalid given name")
        return status

    @staticmethod
    def __is_valid_content(content: str) -> bool:
        status = True if not GenericUtils.str_is_empty(content) else False
        if not status:
            LoggerUtils.error_log("Invalid given content")
        return status

    @staticmethod
    def __get_content_data_with_name(name: str) -> str|None:
        if SHELL_UTILS.is_powershell:
            return f"function {name}"
        elif SHELL_UTILS.is_shell([EShell.BASH, EShell.KSH, EShell.ZSH]):
            return f"function {name}"
        elif SHELL_UTILS.is_fish:
            return f"function {name};"
        return None

    @staticmethod
    def __get_alias_name_from_line(data: str) -> str:
        name: str = data.split("function")[1].split("{")[0].strip() if SHELL_UTILS.is_powershell else data.split("alias")[1].split("=")[0].strip()
        return name if len(name) > 1 else data

    def __delete_process(self, name: str, is_system: bool = False, verbose: bool = False):
        content_data: str | None = AliasProcessor.__get_content_data_with_name(name)
        if content_data is not None:
            FileUtils.delete_line_on_file(self.__alias_file, content_data)
            if is_system:
                content_data = f"Remove-Item Alias:{name} -ErrorAction SilentlyContinue"
                FileUtils.write_file(self.__alias_file, content_data, WriteFileOptions(mode="a"))
            if verbose:
                LoggerUtils.ok_log(f"Deleted alias with name: {name}.")
        else:
            if verbose:
                LoggerUtils.error_log(f"Failed on delete alias with name: {name}")

    def get_all(self, filter_name: str|None) -> list[str]:
        alias_list: list[str] = []
        if self.__is_valid():
            alias_file_data = FileUtils.read_file(self.__alias_file)
            if alias_file_data is not None and len(alias_file_data) > 0:
                for alias_line in alias_file_data.splitlines():
                    alias_name = AliasProcessor.__get_alias_name_from_line(alias_line)
                    if len(alias_name) > 0:
                        alias_list.append(alias_name)
        if filter_name is not None and len(filter_name) > 0:
            alias_list = [alias for alias in alias_list if filter_name in alias]
        return alias_list

    def add(self, name: str, content: str):
        if self.__is_valid() and AliasProcessor.__is_valid_name(name) and AliasProcessor.__is_valid_content(content):
            status: bool = False
            content_data = AliasProcessor.__get_content_data_with_name(name)

            if content_data is not None:
                self.__delete_process(name)
                if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
                    content_data = content_data + " { " + f"{content} {CONST.POWERSHELL_ALL_ARGS_VAR_STR}" + " }"
                    status = FileUtils.write_file(self.__alias_file, content_data, WriteFileOptions(mode="a"))
                    if status and SYSTEM_UTILS.is_windows:
                        script_cmd_file = FileUtils.resolve_path(f"{self.__scripts_dir}/{EShell.CMD.value}/{name}.cmd")
                        cmd_content = self.__windows_cmd_content_template.format(name, CONST.CMD_ALL_ARGS_VAR_STR)
                        status = FileUtils.write_file(script_cmd_file, cmd_content)
                elif SHELL_UTILS.is_shell([EShell.BASH, EShell.KSH, EShell.ZSH]):
                    content_data = content_data + " { " + f"{content} \"{CONST.BASH_ZSH_KSH_ALL_ARGS_VAR_STR}\"" + "; }"
                    status = FileUtils.write_file(self.__alias_file, content_data, WriteFileOptions(mode="a"))
                else:
                    content_data = f"{content_data} {content} {CONST.FISH_ALL_ARGS_VAR_STR}; end"
                    status = FileUtils.write_file(self.__alias_file, content_data, WriteFileOptions(mode="a"))
            if status:
                LoggerUtils.ok_log(f"Added alias with name: {name}")
            else:
                LoggerUtils.error_log(f"Failed on added alias with name: {name}")

    def delete(self, name: str, is_system: bool):
        if self.__is_valid() and AliasProcessor.__is_valid_name(name):
            self.__delete_process(name, is_system, True)
