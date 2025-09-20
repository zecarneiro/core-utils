import re
from dataclasses import dataclass

from coreutils.libs.const_lib import SYSTEM_UTILS, SHELL_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.generic_lib import set_file_permission_to_run
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


@dataclass
class ScriptProcessor:
    @property
    def __windows_cmd_content_template(self) -> str:
        return """
@echo off
powershell.exe -NoLogo -ExecutionPolicy Bypass -File "{0}" {1}"""

    @property
    def __valid_shell(self) -> list[EShell]:
        return [EShell.POWERSHELL, EShell.CMD, EShell.BASH, EShell.KSH, EShell.FISH, EShell.ZSH]

    @property
    def __is_valid_shell(self) -> bool:
        status = SHELL_UTILS.is_shell(self.__valid_shell)
        if not status:
            MessageProcessor.show_shell_msg(self.__valid_shell, self.__class__.__name__)
        return status

    @property
    def __scripts_dir(self) -> str:
        if self.__is_valid_shell:
            directory = DirsLib.get_coreutils_shell_script()
            FileUtils.create_dir(directory)
            return directory
        return ""

    @property
    def __scripts_cmd_dir(self) -> str:
        if SYSTEM_UTILS.is_windows and self.__is_valid_shell:
            directory = FileUtils.resolve_path(f"{self.__scripts_dir}/{EShell.CMD.value}")
            FileUtils.create_dir(directory)
            return directory
        return ""

    @property
    def __is_scripts_dir_created(self) -> bool:
        message = "Failed to create dir"
        status = FileUtils.is_dir(self.__scripts_dir)
        if not status:
            LoggerUtils.error_log(f"{message}: {self.__scripts_dir}")
            return False
        if SYSTEM_UTILS.is_windows:
            status = FileUtils.is_dir(self.__scripts_cmd_dir)
            if not status:
                LoggerUtils.error_log(f"{message}: {self.__scripts_cmd_dir}")
                return False
        return True

    def __is_validate_name(self, name: str) -> bool:
        if not bool(re.fullmatch(r"[A-Za-z0-9_-]+", name)):
            LoggerUtils.error_log("Name accept only A-Z, a-z, 0-9, - and _")
            return False
        return True

    def __is_validate_content(self, content: str) -> bool:
        if GenericUtils.str_is_empty(content):
            LoggerUtils.error_log("Invalid given content")
            return False
        return True

    def __is_valid(self) -> bool:
        if not self.__is_valid_shell or not self.__is_scripts_dir_created:
            return False
        return True

    def __get_bash_bin_script_tag_line(self, data: str) -> str:
        if SHELL_UTILS.is_bash:
            return f"#!/usr/bin/env bash{CONST.EOF}{data}"
        return data

    def get_all(self, filter_name: str|None) -> list[str]:
        script_list: list[str] = []
        if self.__is_valid():
            for file in FileUtils.get_list_files_on_folder(self.__scripts_dir):
                script_list.append(FileUtils.filename_without_ext(FileUtils.basename(file)))
            if not SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
                script_list = script_list + ["pipc", "sudoexe", "appimage-manager"]
        if filter_name is not None and len(filter_name) > 0:
            script_list = [script for script in script_list if filter_name in script]
        return script_list

    def uninstall(self, name: str):
        if self.__is_valid() and self.__is_validate_name(name):
            file_to_delete_list: list[str] = []
            if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
                file_to_delete_list.append(FileUtils.resolve_path(f"{self.__scripts_dir}/{name}.ps1"))
                file_to_delete_list.append(FileUtils.resolve_path(f"{self.__scripts_cmd_dir}/{name}.cmd"))
            elif SHELL_UTILS.is_shell([EShell.BASH]):
                file_to_delete_list.append(FileUtils.resolve_path(f"{self.__scripts_dir}/{name}"))
            for file_to_delete in file_to_delete_list:
                if FileUtils.is_file(file_to_delete):
                    if FileUtils.delete_file(file_to_delete):
                        LoggerUtils.ok_log(f"Deleted file: {file_to_delete}")
                    else:
                        LoggerUtils.error_log(f"Failed to delete file: {file_to_delete}")

    def install(self, name: str, content: str, include_shell_tag: bool = False):
        if self.__is_valid() and self.__is_validate_name(name) and self.__is_validate_content(content):
            success = True
            if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
                script_file = FileUtils.resolve_path(f"{self.__scripts_dir}/{name}.ps1")
                success = FileUtils.write_file(script_file, content)
                if success and SYSTEM_UTILS.is_windows:
                    script_cmd_file = FileUtils.resolve_path(f"{self.__scripts_dir}/{EShell.CMD.value}/{name}.cmd")
                    cmd_content = self.__windows_cmd_content_template.format(script_file, CONST.CMD_ALL_ARGS_VAR_STR)
                    success = FileUtils.write_file(script_cmd_file, cmd_content)
            elif SHELL_UTILS.is_bash:
                if include_shell_tag:
                    content = self.__get_bash_bin_script_tag_line(content)
                script_file = FileUtils.resolve_path(f"{self.__scripts_dir}/{name}")
                success = FileUtils.write_file(script_file, content)
                set_file_permission_to_run(script_file)
            if success:
                LoggerUtils.ok_log(f"Installed script with name: {name}")
            else:
                LoggerUtils.error_log(f"Can not install this script: {name}")

    def install_from_file(self, script_file: str):
        if GenericUtils.str_is_empty(script_file) or not FileUtils.is_file(script_file):
            LoggerUtils.error_log(f"Can not install this script: {script_file}. File not found")
        else:
            name = FileUtils.filename_without_ext(FileUtils.basename(script_file))
            content = FileUtils.read_file(script_file)
            content = content if content is not None else ""
            self.install(name, content)
