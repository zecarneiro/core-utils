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
    name: str
    content: str
    install_file: str
    is_valid_shell: bool = True
    valid_shell: list[EShell] | None = None
    scripts_dir: str = DirsLib.get_coreutils_shell_script()
    windows_cmd_content_template = """
@echo off
powershell -NoProfile -ExecutionPolicy Bypass -File "{0}" {1}"""

    def __post_init__(self):
        self.valid_shell = [EShell.POWERSHELL, EShell.BASH]
        self.is_valid_shell = SHELL_UTILS.is_shell(self.valid_shell)

    def is_valid(self, skip_content_validation: bool = False) -> bool:
        if not self.is_valid_shell:
            MessageProcessor.show_shell_msg(self.valid_shell)
            return False
        elif GenericUtils.str_is_empty(self.name) and GenericUtils.str_is_empty(self.install_file):
            return False
        elif not GenericUtils.str_is_empty(self.name):
            if not bool(re.fullmatch(r"[A-Za-z0-9_-]+", self.name)):
                LoggerUtils.error_log("Name accept only A-Z, a-z, 0-9, - and _")
                return False
            if not skip_content_validation and GenericUtils.str_is_empty(self.content):
                return False
        elif not GenericUtils.str_is_empty(self.install_file) and not FileUtils.is_file(self.install_file):
            LoggerUtils.error_log(f"Can not install this script: {self.install_file}. File not found")
            return False
        return True

    def uninstall(self):
        if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
            file_to_delete_list: list[str] = [
                FileUtils.resolve_path(f"{self.scripts_dir}/{self.name}.ps1"),
                FileUtils.resolve_path(f"{self.scripts_dir}/{EShell.CMD.value}/{self.name}.cmd")
            ]
            for file_to_delete in file_to_delete_list:
                if FileUtils.is_file(file_to_delete):
                    if FileUtils.delete_file(file_to_delete):
                        LoggerUtils.ok_log(f"Deleted file: {file_to_delete}")
                    else:
                        LoggerUtils.error_log(f"Failed to delete file: {file_to_delete}")
        elif SHELL_UTILS.is_shell([EShell.BASH]):
            file_to_delete = FileUtils.resolve_path(f"{self.scripts_dir}/{self.name}")
            if FileUtils.is_file(file_to_delete):
                if FileUtils.delete_file(file_to_delete):
                    LoggerUtils.ok_log(f"Deleted file: {file_to_delete}")
                else:
                    LoggerUtils.error_log(f"Failed to delete file: {file_to_delete}")

    def install_from_file(self):
        if not GenericUtils.str_is_empty(self.install_file) and FileUtils.is_file(self.install_file):
            self.name = FileUtils.filename_without_ext(FileUtils.basename(self.install_file))
            content = FileUtils.read_file(self.install_file)
            self.content = content if content is not None else ""
            self.install_file = ""
            self.install()

    def install(self):
        success_install = True
        if SHELL_UTILS.is_shell([EShell.POWERSHELL, EShell.CMD]):
            script_powershell_file = FileUtils.resolve_path(f"{self.scripts_dir}/{self.name}.ps1")
            script_cmd_file = FileUtils.resolve_path(f"{self.scripts_dir}/{EShell.CMD.value}/{self.name}.bat")
            cmd_content = self.windows_cmd_content_template.format(script_powershell_file, CONST.CMD_ALL_ARGS_VAR_STR)
            success_install = FileUtils.write_file(script_powershell_file, self.content)
            if success_install and SYSTEM_UTILS.is_windows:
                success_install = FileUtils.write_file(script_cmd_file, cmd_content)
        elif SHELL_UTILS.is_bash:
            bash_env_import = "#!/usr/bin/env bash"
            bash_import = "#!/usr/bin/bash"
            if bash_env_import not in self.content and bash_import not in self.content:
                self.content = f"{bash_env_import}{CONST.EOF}{self.content}"
            script_bash_file = FileUtils.resolve_path(f"{self.scripts_dir}/{self.name}")
            success_install = FileUtils.write_file(script_bash_file, self.content)
            set_file_permission_to_run(script_bash_file)
        if success_install:
            LoggerUtils.ok_log(f"Installed script with name: {self.name}")
        else:
            LoggerUtils.error_log(f"Can not install this script: {self.name}")

    def get_all(self) -> list[str]:
        script_list: list[str] = []
        for file in FileUtils.get_list_files_on_folder(self.scripts_dir):
            script_list.append(FileUtils.filename_without_ext(FileUtils.basename(file)))
        return script_list