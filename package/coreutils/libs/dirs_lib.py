import importlib.resources as res
import coreutils
from coreutils.libs.const_lib import SYSTEM_UTILS, SHELL_UTILS
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from pathlib import Path
from coreutils.libs.pythonutils.file_utils import FileUtils


class DirsLib:
    @staticmethod
    def __create_dir(directory: str) -> str:
        if not FileUtils.is_dir(directory):
            FileUtils.create_dir(directory)
        return directory

    @staticmethod
    def get_config() -> str:
        return DirsLib.__create_dir(FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}/.config"))

    @staticmethod
    def get_user_local() -> str:
        return DirsLib.__create_dir(FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}/.local"))

    @staticmethod
    def get_user_opt() -> str:
        return DirsLib.__create_dir(FileUtils.resolve_path(f"{DirsLib.get_user_local()}/opt"))

    @staticmethod
    def get_user_bin() -> str:
        return DirsLib.__create_dir(FileUtils.resolve_path(f"{DirsLib.get_user_local()}/bin"))

    @staticmethod
    def get_user_startup() -> str|None:
        if SYSTEM_UTILS.is_windows:
            return DirsLib.__create_dir(FileUtils.resolve_path(f"{SYSTEM_UTILS.home_dir}\\Start Menu\\Programs\\Startup"))
        elif SYSTEM_UTILS.is_linux:
            return DirsLib.__create_dir(FileUtils.resolve_path(f"{DirsLib.get_config()}/autostart"))
        return None

    @staticmethod
    def get_user_temp() -> str | None:
        if SYSTEM_UTILS.is_windows:
            return DirsLib.__create_dir(FileUtils.resolve_path(SYSTEM_UTILS.temp_dir))
        elif SYSTEM_UTILS.is_linux:
            return DirsLib.__create_dir(FileUtils.resolve_path(f"{DirsLib.get_user_local()}/tmp"))
        return None

    @staticmethod
    def get_temp() -> str | None:
        if SYSTEM_UTILS.is_windows:
            return DirsLib.get_user_temp()
        elif SYSTEM_UTILS.is_linux:
            return FileUtils.resolve_path(SYSTEM_UTILS.temp_dir)
        return None

    @staticmethod
    def get_coreutils_local_dir() -> str:
        return DirsLib.__create_dir(FileUtils.resolve_path(f"{DirsLib.get_user_local()}/coreutils"))

    @staticmethod
    def get_coreutils_shell_function() -> str:
        current_shell = SHELL_UTILS.current_shell
        if current_shell == EShell.UNKNOWN:
            return ""
        shell_dir_name = current_shell.value if current_shell != EShell.CMD else EShell.POWERSHELL.value
        directory = FileUtils.resolve_path(f"{DirsLib.get_coreutils_local_dir()}/functions/{shell_dir_name}")
        if not FileUtils.is_dir(directory):
            FileUtils.create_dir(directory)
            if SYSTEM_UTILS.is_windows:
                FileUtils.create_dir(FileUtils.resolve_path(f"{directory}/{EShell.CMD.value}"))
        return directory

    @staticmethod
    def get_coreutils_shell_script() -> str:
        current_shell = SHELL_UTILS.current_shell
        if current_shell == EShell.UNKNOWN:
            return ""
        shell_dir_name = current_shell.value if current_shell != EShell.CMD else EShell.POWERSHELL.value
        directory = FileUtils.resolve_path(f"{DirsLib.get_coreutils_local_dir()}/scripts/{shell_dir_name}")
        if not FileUtils.is_dir(directory):
            FileUtils.create_dir(directory)
            if SYSTEM_UTILS.is_windows:
                FileUtils.create_dir(FileUtils.resolve_path(f"{directory}/{EShell.CMD.value}"))
        return directory

    @staticmethod
    def get_coreutils_shell_alias() -> str:
        current_shell = SHELL_UTILS.current_shell
        if current_shell == EShell.UNKNOWN:
            return ""
        directory = FileUtils.resolve_path(f"{DirsLib.get_coreutils_local_dir()}/alias")
        if not FileUtils.is_dir(directory):
            FileUtils.create_dir(directory)
            if SYSTEM_UTILS.is_windows:
                FileUtils.create_dir(FileUtils.resolve_path(f"{directory}/{EShell.CMD.value}"))
        return directory

    @staticmethod
    def get_coreutils_config_dir() -> str:
        return DirsLib.__create_dir(FileUtils.resolve_path(f"{DirsLib.get_config()}/coreutils"))

    @staticmethod
    def get_coreutils_dir() -> Path:
        return Path(coreutils.__file__).resolve().parent

    @staticmethod
    def get_resource_dir() -> Path:
        local_resources = DirsLib.get_coreutils_dir() / "resources"
        if local_resources.exists():
            # Running from source or editable install
            return local_resources
        # Installed as wheel (zipped) → use importlib.resources
        with res.as_file(res.files("coreutils") / "resources") as path:
            return path

    @staticmethod
    def get_resource_shell_script_dir() -> str:
        return FileUtils.resolve_path(f"{DirsLib.get_resource_dir()}/shell-scripts")

    @staticmethod
    def get_resource_shell_script_libs_dir() -> str:
        return FileUtils.resolve_path(f"{DirsLib.get_resource_shell_script_dir()}/libs")
