from dataclasses import dataclass

from coreutils.libs.const_lib import SYSTEM_UTILS, SHELL_UTILS
from coreutils.libs.pythonutils.enums.platform_enum import EPlatform
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.logger_utils import LoggerUtils


@dataclass
class MessageProcessor:
    @staticmethod
    def show_platform_msg(platform_list: list[EPlatform]):
        base_msg = "Running __OS__-specific command only"
        msg: str | None = None
        count_platform = 0
        if platform_list is not None:
            if EPlatform.WINDOWS in platform_list and not SYSTEM_UTILS.is_windows:
                count_platform = count_platform + 1
            if EPlatform.LINUX in platform_list and not SYSTEM_UTILS.is_linux:
                count_platform = count_platform + 1
            if count_platform >= len(platform_list):
                valid_platform = "["
                for platform in platform_list:
                    if len(valid_platform) > 2:
                        valid_platform = f"{valid_platform}, {platform.name}"
                    else:
                        valid_platform = f"{valid_platform}{platform.name}"
                valid_platform = f"{valid_platform}]"
                msg = base_msg.replace("__OS__", valid_platform)
        if msg is not None:
            LoggerUtils.warn_log(msg)

    @staticmethod
    def show_shell_msg(shell_list: list[EShell]):
        base_msg = "Running __SHELL__-specific command only"
        msg: str | None = None
        if not SHELL_UTILS.is_shell(shell_list):
            valid_shell = "["
            for shell in shell_list:
                if len(valid_shell) > 2:
                    valid_shell = f"{valid_shell}, {shell.value}"
                else:
                    valid_shell = f"{valid_shell}{shell.value}"
            valid_shell = f"{valid_shell}]"
            msg = base_msg.replace("__SHELL__", valid_shell)
        if msg is not None:
            LoggerUtils.warn_log(msg)