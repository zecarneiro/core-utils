import importlib.resources as res
from pathlib import Path

import coreutils
from coreutils.libs.const_lib import SHELL_UTILS, SYSTEM_UTILS
from coreutils.libs.pythonutils.enums.shell_enum import EShell
from coreutils.libs.pythonutils.file_utils import FileUtils


class DirsLib:
    @staticmethod
    def __create_dir(directory: str) -> str:
        if not FileUtils.is_dir(directory):
            FileUtils.create_dir(directory)
        return directory

    @staticmethod
    def get_user_bin() -> str:
        return DirsLib.__create_dir(
            FileUtils.resolve_path(f"{DirsLib.get_user_local()}/bin")
        )

    @staticmethod
    def get_coreutils_local_dir() -> str:
        return DirsLib.__create_dir(
            FileUtils.resolve_path(f"{DirsLib.get_user_local()}/coreutils")
        )

    @staticmethod
    def get_coreutils_shell_function() -> str:
        current_shell = SHELL_UTILS.current_shell
        if current_shell == EShell.UNKNOWN:
            return ""
        shell_dir_name = (
            current_shell.value
            if current_shell != EShell.CMD
            else EShell.POWERSHELL.value
        )
        directory = FileUtils.resolve_path(
            f"{DirsLib.get_coreutils_local_dir()}/functions/{shell_dir_name}"
        )
        if not FileUtils.is_dir(directory):
            FileUtils.create_dir(directory)
            if SYSTEM_UTILS.is_windows:
                FileUtils.create_dir(
                    FileUtils.resolve_path(f"{directory}/{EShell.CMD.value}")
                )
        return directory

    @staticmethod
    def get_coreutils_config_dir() -> str:
        return DirsLib.__create_dir(
            FileUtils.resolve_path(f"{DirsLib.get_config()}/coreutils")
        )

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
    def get_resource_shell_script_dir(file: str | None = None) -> str:
        if file is not None and len(file) > 0:
            return FileUtils.resolve_path(
                f"{DirsLib.get_resource_dir()}/shell-scripts/{file}"
            )
        return FileUtils.resolve_path(f"{DirsLib.get_resource_dir()}/shell-scripts")

    @staticmethod
    def get_resource_shell_script_libs_dir(lib: str | None = None) -> str:
        directory = DirsLib.get_resource_shell_script_dir("libs")
        return (
            directory
            if lib is None or len(lib) == 0
            else FileUtils.resolve_path(f"{directory}/{lib}")
        )

    @staticmethod
    def get_resource_shell_script_apps_dir(app: str | None = None) -> str:
        directory = DirsLib.get_resource_shell_script_dir("apps")
        return (
            directory
            if app is None or len(app) == 0
            else FileUtils.resolve_path(f"{directory}/{app}")
        )
