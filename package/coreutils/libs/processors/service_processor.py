from dataclasses import dataclass

from coreutils.libs.const_lib import SYSTEM_UTILS, CONSOLE_UTILS
from coreutils.libs.dirs_lib import DirsLib
from coreutils.libs.processors.message_processor import MessageProcessor
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.platform_enum import EPlatform
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.libs.pythonutils.system_utils import SystemUtils


@dataclass
class ServiceProcessor:
    name: str
    command: str
    valid_platform: list[EPlatform]|None = None
    description: str = ""
    service_dir = FileUtils.resolve_path(f"{DirsLib.get_coreutils_local_dir()}/services")
    is_service_dir_created = False
    is_valid_platform = True
    windows_data = """
@echo off
echo {0}
{1}
"""
    linux_data = """
[Unit]
Description={0}

[Service]
ExecStart={1}
Restart=always
StandardOutput=append:/var/log/{2}.log
StandardError=append:/var/log/{2}.log

[Install]
WantedBy=multi-user.target
"""

    def __post_init__(self):
        self.valid_platform = [EPlatform.LINUX, EPlatform.WINDOWS]
        self.is_valid_platform = SYSTEM_UTILS.is_platform(self.valid_platform)
        if self.is_valid_platform:
            self.is_service_dir_created = FileUtils.create_dir(self.service_dir)
        else:
            MessageProcessor.show_platform_msg(self.valid_platform)

    def is_valid(self, skip_command: bool = False) -> bool:
        if not self.is_valid_platform:
            return False
        elif not self.is_service_dir_created:
            LoggerUtils.error_log(f"Create service dir failed: {self.service_dir}")
            return False
        elif GenericUtils.str_is_empty(self.name):
            LoggerUtils.error_log(f"Invalid given name: {self.name}")
            return False
        elif not skip_command and GenericUtils.str_is_empty(self.command):
            LoggerUtils.error_log(f"Invalid given command: {self.command}")
            return False
        elif not SystemUtils.is_admin():
            LoggerUtils.error_log(f"Please, run as admin to continue...")
            return False
        return True

    def uninstall(self):
        success_uninstall = True
        if SYSTEM_UTILS.is_windows:
            bat_file = FileUtils.resolve_path(f"{self.service_dir}/{self.name}.bat")
            cmd_list = [
                f"sudo sc.exe stop \"{self.name}\"",
                f"sudo sc.exe delete \"{self.name}\""
            ]
            for cmd in cmd_list:
                CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL, verbose=True))
            success_uninstall = FileUtils.delete_file(bat_file)
        elif SYSTEM_UTILS.is_linux:
            service_file = FileUtils.resolve_path(f"/etc/systemd/system/{self.name}.service")
            cmd_list = [
                f"sudo systemctl stop {self.name}",
            ]
            for cmd in cmd_list:
                CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, verbose=True))
            success_uninstall = FileUtils.delete_file(service_file)
        if success_uninstall:
            LoggerUtils.ok_log(f"Uninstall service with name: {self.name}")
        else:
            LoggerUtils.error_log(f"Can not uninstall this service: {self.name}")

    def install(self):
        success_install = True
        if SYSTEM_UTILS.is_windows:
            bat_file = FileUtils.resolve_path(f"{self.service_dir}/{self.name}.bat")
            success_install = FileUtils.write_file(bat_file, self.windows_data.format(self.description, self.command))
            if success_install:
                cmd_list = [
                    f"sudo sc.exe create \"{self.name}\" binPath=\"{bat_file}\" type=own start=auto",
                    f"sudo sc.exe start \"{self.name}\""
                ]
                for cmd in cmd_list:
                    CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, shell=EShell.POWERSHELL, verbose=True))
        elif SYSTEM_UTILS.is_linux:
            service_file = FileUtils.resolve_path(f"/etc/systemd/system/{self.name}.service")
            success_install = FileUtils.write_file(service_file, self.linux_data.format(self.description, self.command, self.name))
            if success_install:
                cmd_list = [
                    f"sudo systemctl start {self.name}",
                    f"sudo systemctl enable {self.name}"
                ]
                for cmd in cmd_list:
                    CONSOLE_UTILS.exec_real_time(CommandInfo(command=cmd, verbose=True))
        if success_install:
            LoggerUtils.ok_log(f"Installed service with name: {self.name}")
        else:
            LoggerUtils.error_log(f"Can not install this service: {self.name}")
