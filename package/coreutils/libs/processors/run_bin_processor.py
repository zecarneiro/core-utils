from dataclasses import dataclass

from coreutils.libs.const_lib import CONSOLE_UTILS
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


@dataclass
class RunBinProcessor:
    @property
    def __pwsh_start_process(self) -> str:
        return "Start-Process '{0}' -Wait"

    @property
    def __pwsh_add_appx_package(self) -> str:
        return "Add-AppxPackage -Path \"{0}\""

    def __run_bin(self, binary: str, verbose: bool = True):
        if FileUtils.is_file(binary):
            file_ext = FileUtils.file_extension(binary)
            match file_ext:
                case ".exe" | ".msi":
                    cmd = CommandInfo(command=self.__pwsh_start_process.format(binary), shell=EShell.POWERSHELL)
                    CONSOLE_UTILS.exec_real_time(cmd)
                case ".msixbundle":
                    cmd = CommandInfo(command=self.__pwsh_add_appx_package.format(binary), shell=EShell.POWERSHELL)
                    CONSOLE_UTILS.exec_real_time(cmd)
                case _:
                    if verbose:
                        LoggerUtils.error_log(f"Can not run this bin file: {binary}")
                        LoggerUtils.info_log("Only accept file with ext: .exe|.msi|.msixbundle")
        else:
            LoggerUtils.error_log("Invalid given bin file")

    def __run_dir(self, directory: str):
        if FileUtils.is_dir(directory):
            files = FileUtils.get_list_files_on_folder(directory)
            for file in files:
                self.__run_bin(file, False)
            LoggerUtils.info_log(f"Execution of all files on '{directory}' it's done.")

    def start(self, binary_or_dir: str):
        if not GenericUtils.str_is_empty(binary_or_dir):
            if FileUtils.is_file(binary_or_dir):
                self.__run_bin(binary_or_dir)
            elif FileUtils.is_dir(binary_or_dir):
                self.__run_dir(binary_or_dir)

