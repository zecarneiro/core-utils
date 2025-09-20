from dataclasses import dataclass, field

from coreutils.libs.const_lib import CONSOLE_UTILS, SHELL_UTILS
from coreutils.libs.pythonutils.console_utils import ConsoleUtils
from coreutils.libs.pythonutils.const_utils import CONST
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.generic_utils import GenericUtils


@dataclass
class WhichProcessor:
    command: str
    __result: str = field(default="", init=False)
    __function_key: str = field(default="FUNCTION", init=False)
    __alias_key: str = field(default="ALIAS", init=False)

    def __has_command(self) -> bool:
        return not GenericUtils.str_is_empty(self.__result)

    def __run_command_info(self, command_info: CommandInfo|None, type_command: str):
        if command_info is not None:
            res_exec = CONSOLE_UTILS.exec(command_info)
            if len(res_exec.stdout) > 0:
                self.__result = f"{type_command}: {self.command}"

    def __get_command_info(self, type_command: str) -> CommandInfo|None:
        if SHELL_UTILS.is_powershell:
            if type_command == self.__function_key:
                return CommandInfo(command=f"Get-ChildItem Function:{self.command} -ErrorAction SilentlyContinue")
            elif type_command == self.__alias_key:
                return CommandInfo(command=f"Get-Alias {self.command} -ErrorAction SilentlyContinue")
        elif SHELL_UTILS.is_bash:
            if type_command == self.__function_key:
                cmd = f"bash -i -c 'declare -F {self.command} >/dev/null 2>&1 && echo existe'"
                return CommandInfo(command=cmd, shell=EShell.BASH, use_shell=True)
            elif type_command == self.__alias_key:
                cmd = f"bash -i -c 'alias {self.command} >/dev/null 2>&1 && echo existe'"
                return CommandInfo(command=cmd, shell=EShell.BASH, use_shell=True)
        elif SHELL_UTILS.is_fish:
            if type_command == self.__function_key:
                return CommandInfo(command=f"functions -q {self.command}; and echo \"exists\"", shell=EShell.FISH, use_shell=True)
            elif type_command == self.__alias_key:
                cmd = f"alias | grep -q '^alias {self.command} '; and echo exists"
                return CommandInfo(command=cmd, shell=EShell.FISH, use_shell=True)
        elif SHELL_UTILS.is_ksh or SHELL_UTILS.is_zsh:
            prefix_cmd = "ksh -i -c" if SHELL_UTILS.is_ksh else "zsh -i -c"
            if type_command == self.__function_key:
                return CommandInfo(command=f"{prefix_cmd} 'typeset -f {self.command} >/dev/null 2>&1 && echo exists'", use_shell=True)
            elif type_command == self.__alias_key:
                cmd = f"{prefix_cmd} 'alias | grep -q \"^{self.command}=\" && echo exists'"
                return CommandInfo(command=cmd, use_shell=True)
        return None

    def __find_app_file_path(self) -> bool:
        self.__result = ConsoleUtils.which(self.command)
        status = self.__has_command()
        if status:
            self.__result = f"APP/CMD: {self.__result}"
        return status

    def __find_function(self) -> bool:
        command_info = self.__get_command_info(self.__function_key)
        self.__run_command_info(command_info, self.__function_key)
        return self.__has_command()

    def __find_alias(self) -> bool:
        command_info = self.__get_command_info(self.__alias_key)
        self.__run_command_info(command_info, self.__alias_key)
        return self.__has_command()

    def find_command(self):
        if self.__find_app_file_path() or self.__find_alias() or self.__find_function():
            return self.__result
        return CONST.UNKNOWN
