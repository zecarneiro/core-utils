import argparse

from coreutils.libs.const_lib import CONSOLE_UTILS
from coreutils.libs.pythonutils.entities.command_info import CommandInfo
from coreutils.libs.pythonutils.entities.github_repo_release_options import GitHubRepoReleaseOption
from coreutils.libs.pythonutils.file_utils import FileUtils
from coreutils.libs.pythonutils.generic_utils import GenericUtils
from coreutils.libs.pythonutils.logger_utils import LoggerUtils
from coreutils.systemfunctions import script_processor

__GIT_BASE_CMD: str = "git"

def __get_git_command_info() -> CommandInfo:
    return CommandInfo(command=__GIT_BASE_CMD, verbose=True)

def git_dependencies_apps() -> list[str]:
    apps: list[str] = [__GIT_BASE_CMD]
    return apps

def process_post_install_for_git_function_file():
    script_processor(["-i", "-n", "git-undo-last-commit", "-c", f"{__GIT_BASE_CMD} reset --soft HEAD~1"])
    script_processor(["-i", "-n", "git-stage-all", "-c", f"{__GIT_BASE_CMD} add ."])
    script_processor(["-i", "-n", "git-status", "-c", f"{__GIT_BASE_CMD} status"])
    script_processor(["-i", "-n", "git-cherry-pick-master-continue", "-c", f"{__GIT_BASE_CMD} cherry-pick --continue"])

def git_reset_hard_origin():
    command_info = __get_git_command_info()
    # Get branch name
    command_info.args = ["branch", "--show-current"]
    response = CONSOLE_UTILS.exec(command_info)
    if len(response.stderr) > 0:
        LoggerUtils.error_log(response.stderr)
    else:
        current_branch_name = response.stdout
        command_info.args = ["reset", "--hard", f"origin/{current_branch_name}"]
        CONSOLE_UTILS.exec_real_time(command_info)

def git_reset_file():
    branch: str = "origin/master"
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", metavar="FILEPATH", type=str, required=True)
    parser.add_argument("-b", "--branch", metavar="BRANCH", type=str, default=branch, help="Name of branch you want to file reset to. (DEFAULT: origin/master)")
    args = parser.parse_args()
    file: str = args.file
    if FileUtils.file_exist(file):
        command_info = __get_git_command_info()
        command_info.args = ["checkout", branch, f"'{file}'"]
        CONSOLE_UTILS.exec_real_time(command_info)
    else:
        LoggerUtils.error_log(f"Invalid given file: {file}")

def git_repo_backup():
    parser = argparse.ArgumentParser()
    parser.add_argument("-u", "--url", metavar="URL", type=str, required=True)
    args = parser.parse_args()
    command_info = __get_git_command_info()
    command_info.args = ["clone", "--mirror", f"'{args.url}'"]
    CONSOLE_UTILS.exec_real_time(command_info)

def git_repo_restore_backup():
    parser = argparse.ArgumentParser()
    parser.add_argument("-u", "--url", metavar="URL", type=str, required=True)
    args = parser.parse_args()
    command_info = __get_git_command_info()
    command_info.args = ["push", "--mirror", f"'{args.url}'"]
    CONSOLE_UTILS.exec_real_time(command_info)

def git_move_submodule():
    parser = argparse.ArgumentParser()
    parser.add_argument("-o", "--old", metavar="PATH", type=str, required=True, help="Old path")
    parser.add_argument("-n", "--new", metavar="PATH", type=str, required=True, help="New path")
    args = parser.parse_args()
    old: str = args.old
    new: str = args.new
    new_parent_dir: str = FileUtils.dirname(new)
    if FileUtils.file_exist(new_parent_dir):
        LoggerUtils.info_log(f"Directory already exists: {new_parent_dir}")
    else:
        FileUtils.create_dir(new_parent_dir)
    command_info = __get_git_command_info()
    command_info.args = ["mv", f"'{old}'", f"'{new}'"]
    CONSOLE_UTILS.exec_real_time(command_info)

def git_add_script_perm():
    parser = argparse.ArgumentParser()
    parser.add_argument("-s", "--script", metavar="SCRIPT", type=str, required=True, help="Script to add permission")
    args = parser.parse_args()
    script_filename = FileUtils.basename(args.script)
    command_info = __get_git_command_info()
    command_info.args = ["update-index", "--chmod=+x", f"'{args.script}'"]
    CONSOLE_UTILS.exec_real_time(command_info)
    command_info = __get_git_command_info()
    command_info.args = ["ls-files", "--stage"]
    res = CONSOLE_UTILS.exec(command_info)
    if len(res.stderr) > 0:
        LoggerUtils.error_log(res.stderr)
    else:
        print(GenericUtils.list_to_str(GenericUtils.grep_e(res.stdout, script_filename, False)))

def git_cherrypick_master():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--commit", metavar="COMMIT", nargs="?", required=True)
    args = parser.parse_args()
    command_info = __get_git_command_info()
    command_info.args = ["cherry-pick", "-m", "1", f"'{GenericUtils.list_to_str(args.commit)}'"]
    CONSOLE_UTILS.exec_real_time(command_info)

def gitclone():
    parser = argparse.ArgumentParser()
    parser.add_argument("-u", "--url", metavar="URL", type=str, required=True)
    args = parser.parse_args()
    command_info = __get_git_command_info()
    command_info.args = ["clone", f"'{args.url}'"]
    CONSOLE_UTILS.exec_real_time(command_info)

def github_change_url():
    parser = argparse.ArgumentParser()
    parser.add_argument("-u", "--username", metavar="USERNAME", type=str, required=True, help="Github USERNAME")
    parser.add_argument("-t", "--token", metavar="TOKEN", type=str, required=True, help="Github TOKEN")
    parser.add_argument("-r", "--repository", metavar="REPO", type=str, required=True, help="Github new REPOSITORY_NAME (ex: REPOSITORY_NAME = AAA/bbb.gi or github.com/REPOSITORY_NAME)")
    args = parser.parse_args()
    username: str = args.username
    token: str = args.token
    url: str = args.url
    full_url: str = f"https://${username}:${token}@github.com/{url}"
    LoggerUtils.info_log(f"Will be set new github URL: {full_url}")
    command_info = __get_git_command_info()
    command_info.args = ["remote", "set-url", "origin", f"'{full_url}'"]
    CONSOLE_UTILS.exec_real_time(command_info)

def git_set_config():
    git_dir_or_file_name = ".git"
    local_cmd_list: list[str] = [
        "config core.autocrlf input",
        "config core.fileMode false",
        "config core.logAllRefUpdates true",
        "config core.ignorecase true",
        "config pull.rebase true",
        "config --unset safe.directory",
        "config --add safe.directory '*'",
        "config merge.ff false",
    ]
    global_cmd_list: list[str] = [
        (cmd.replace("config", "config --global") for cmd in local_cmd_list),
        "config --global credential.credentialStore plaintext",
        "git-credential-manager configure",
    ]
    # Set global configs
    for cmd in global_cmd_list:
        command_info = __get_git_command_info()
        command_info.args = [cmd]
        CONSOLE_UTILS.exec_real_time(command_info)
    if FileUtils.file_exist(git_dir_or_file_name):
        LoggerUtils.info_log("Set local configurations")
        for cmd in local_cmd_list:
            command_info = __get_git_command_info()
            command_info.args = [cmd]
            CONSOLE_UTILS.exec_real_time(command_info)

def git_config_user():
    parser = argparse.ArgumentParser()
    parser.add_argument("-u", "--username", metavar="USERNAME", type=str)
    parser.add_argument("-e", "--email", metavar="EMAIL", type=str)
    args = parser.parse_args()
    username: str|None = args.username
    email: str|None = args.email
    global_flag = "--global" if FileUtils.file_exist(".git") else ""
    command_info = __get_git_command_info()
    base_cmd_args = ["config", global_flag]
    already_run = False
    if username is not None and len(username) > 0:
        command_info.args = base_cmd_args + ["user.name", f"'{username}'"]
        CONSOLE_UTILS.exec_real_time(command_info)
        already_run = True
    if email is not None and len(email) > 0:
        command_info.args = base_cmd_args + ["user.email", f"'{email}'"]
        CONSOLE_UTILS.exec_real_time(command_info)
        already_run = True
    if not already_run:
        parser.print_help(None)

def gitcommit():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--commit", metavar="COMMIT", type=str, required=True)
    args = parser.parse_args()
    command_info = __get_git_command_info()
    command_info.args = ["commit", "-m", f"'{args.commit}'"]
    CONSOLE_UTILS.exec_real_time(command_info)

def github_latest_version():
    no_release_found_msg = "No releases found for this repository."
    parser = argparse.ArgumentParser()
    parser.add_argument("-o", "--owner", metavar="OWNER", type=str, required=True, help="Github OWNER")
    parser.add_argument("-r", "--repository", metavar="REPO", type=str, required=True, help="Github REPOSITORY")
    parser.add_argument("-l", "--latest", action="store_true", dest="latest", help="Enable to get latest version")  # store_true = DEFAULT False
    args = parser.parse_args()
    owner: str = args.owner
    repository: str = args.repository
    release = GenericUtils.get_github_repo_release(GitHubRepoReleaseOption(owner=owner, repo=repository, is_latest=args.latest))
    if release is not None:
        already_run = False
        latest_version: str = release.tag_name.replace("v", "") if release.tag_name else ""  # usually the version tag
        if not GenericUtils.str_is_empty(latest_version):
            print(f"Latest version: {latest_version}")
            already_run = True
        if not GenericUtils.str_is_empty(release.zipball_url): # zip archive of this release
            print(f"Download URL: {release.zipball_url}")
            already_run = True
        if not already_run:
            LoggerUtils.warn_log(no_release_found_msg)
    else:
        LoggerUtils.warn_log(no_release_found_msg)
