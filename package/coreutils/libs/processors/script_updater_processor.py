import shutil
from dataclasses import dataclass

from coreutils.libs.const_lib import SHELL_UTILS, CONSOLE_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.generic_lib import set_file_permission_to_run
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


@dataclass
class ScriptUpdaterProcessor:
    @property
    def __valid_shell(self) -> list[EShell]:
        return [EShell.POWERSHELL, EShell.BASH]

    @property
    def __is_valid_shell(self) -> bool:
        status = SHELL_UTILS.is_shell(self.__valid_shell)
        if not status:
            MessageProcessor.show_shell_msg(self.__valid_shell, self.__class__.__name__)
        return status

    @property
    def __shell_dir(self) -> str:
        directory = FileUtils.resolve_path(f"{DirsLib.get_coreutils_local_dir()}/shell-scripts-installers-updaters")
        if self.__is_valid_shell:
            directory = FileUtils.resolve_path(f"{directory}/{SHELL_UTILS.current_shell.value}")
            FileUtils.create_dir(directory)
            return directory
        return ""

    @property
    def __is_shell_dir_created(self) -> bool:
        status = FileUtils.is_dir(self.__shell_dir)
        if not status:
            LoggerUtils.error_log(f"Failed to create dir: {self.__shell_dir}")
        return status

    @property
    def __file_list(self) -> list[str]:
        return FileUtils.get_list_files_on_folder(self.__shell_dir)

    def __is_valid(self) -> bool:
        return self.__is_valid_shell and self.__is_shell_dir_created

    def install(self, file: str):
        if self.__is_valid() and FileUtils.is_file(file):
            file_basename = FileUtils.basename(file)
            dest = FileUtils.resolve_path(f"{self.__shell_dir}/{file_basename}")
            can_install = True if not FileUtils.is_file(dest) else False
            if not can_install and ConsoleUtils.confirm("Script already exists. Continue", False):
                can_install = True
            if can_install:
                shutil.copy2(file, dest)
                set_file_permission_to_run(dest)
                LoggerUtils.ok_log(f"Installed shell script updater: {file_basename}")
        else:
            LoggerUtils.error_log("Installation failed")

    def uninstall(self, name: str):
        if self.__is_valid():
            for file in self.__file_list:
                file_basename = FileUtils.basename(file)
                if name == file_basename or name == FileUtils.filename_without_ext(file_basename):
                    if FileUtils.delete_file(file):
                        LoggerUtils.ok_log(f"Deleted  {file_basename}")
                    else:
                        LoggerUtils.error_log(f"Deleted  {file_basename} was failed")
        else:
            LoggerUtils.error_log("Uninstall failed")

    def get_all(self, filter_name: str|None) -> list[str]:
        script_list: list[str] = []
        if self.__is_valid():
            for file in self.__file_list:
                script_list.append(FileUtils.filename_without_ext(FileUtils.basename(file)))
        if filter_name is not None and len(filter_name) > 0:
            script_list = [script for script in script_list if filter_name in script]
        return script_list

    def run(self, name: str, is_all: bool):
        if self.__is_valid():
            if is_all:
                LoggerUtils.title_log("Process all installed shell script to install/update/uninstall packages")
            for file in self.__file_list:
                file_basename = FileUtils.basename(file)
                can_run = True
                if not is_all:
                    if not (name == FileUtils.filename_without_ext(file_basename) or name == file_basename):
                        can_run = False
                if can_run:
                    LoggerUtils.header_log(f"Running {FileUtils.filename_without_ext(file_basename)}")
                    CONSOLE_UTILS.exec_real_time(CommandInfo(command=f". \"{file}\"", shell=SHELL_UTILS.current_shell))
