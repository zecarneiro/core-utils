from vendor.pythonutils.logger_utils import title_log, header_log

__HEADER_LENGTH = 30


def __install_shell_script():
    header_log("Create all shell scripts", )

def __create_default_dirs():
    header_log("Create all default directories", __HEADER_LENGTH)

def run():
    title_log("Running custom setup logic...")
    # Example logic
    import os
    home = os.path.expanduser("~")
    print(f"Home directory is {home}")