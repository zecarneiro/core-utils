import os
import glob
import configparser
from dataclasses import dataclass

from coreutils.libs.const_lib import CONSOLE_UTILS, SYSTEM_UTILS
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.generic_utils import GenericUtils

@dataclass
class StartAppInfo:
    name: str|None
    exec: str|None
    icon: str|None
    path: str|None

@dataclass
class StartAppsProcessor:
    @property
    def __linux_desktop_dirs(self) -> list[str]:
        return [
            "/usr/share/applications",
            "/var/lib/snapd/desktop/applications",
            "/var/lib/flatpak/exports/share/applications",
            f"{SYSTEM_UTILS.home_dir}/.local/share/applications",
            f"{SYSTEM_UTILS.home_dir}/.local/share/flatpak/exports/share/applications",
        ]

    def __has_filter(self, filter_app: str|None) -> bool:
        return True if not GenericUtils.str_is_empty(filter_app) else False

    def __process_windows(self, filter_app: str, is_case_insensitive: bool):
        command_info = CommandInfo(command="Get-StartApps", shell=EShell.POWERSHELL)
        if self.__has_filter(filter_app):
            res = CONSOLE_UTILS.exec(command_info)
            if res.has_error():
                res.log_error()
            else:
                matches: list[str] = GenericUtils.grep_e(res.stdout, filter_app, False if is_case_insensitive else True)
                for match in matches:
                    print(match)
        else:
            CONSOLE_UTILS.exec_real_time(command_info)

    def __parse_desktop_file(self, path: str) -> StartAppInfo|None:
        """Extract Name, Exec, Icon from a .desktop file."""
        config = configparser.ConfigParser(interpolation=None)
        try:
            config.read(path)
            section = "Desktop Entry"
            name = config.get(section, "Name", fallback=None)
            exec_cmd = config.get(section, "Exec", fallback=None)
            icon = config.get(section, "Icon", fallback=None)
            return StartAppInfo(name=name, exec=exec_cmd, icon=icon, path=path)
        except Exception:
            return None

    def __process_linux(self, filter_app: str, is_case_insensitive: bool):
        desktop_file_list: list[str] = []
        for directory in self.__linux_desktop_dirs:
            if os.path.isdir(directory):
                desktop_file_list.extend(glob.glob(os.path.join(directory, "*.desktop")))

        apps: list[StartAppInfo] = []
        for desktop_file in desktop_file_list:
            data = self.__parse_desktop_file(desktop_file)
            if data and data.name and data.exec:
                apps.append(data)

        for app in apps:
            can_print = True
            if self.__has_filter(filter_app):
                name = app.name if app.name is not None else ""
                matches: list[str] = GenericUtils.grep_e(name, filter_app, False if is_case_insensitive else True)
                if len(matches) == 0:
                    can_print = False
            if can_print:
                print(f"{app.name}")
                print(f"  Exec: {app.exec}")
                print(f"  Icon: {app.icon}")
                print(f"  File: {app.path}")
                print()

    def start(self, filter_app: str, is_case_insensitive: bool):
        if SYSTEM_UTILS.is_windows:
            self.__process_windows(filter_app, is_case_insensitive)
        else:
            self.__process_linux(filter_app, is_case_insensitive)


