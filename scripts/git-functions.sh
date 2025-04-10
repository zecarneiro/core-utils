#!/usr/bin/env bash
# Author: José M. C. Noronha

function gitresethardorigin {
    local current_branch_name="$(git branch --show-current)"
    git reset --hard origin/$current_branch_name
}
function gitresetfile {
    fileName="$1"
    branch="$2"
    if [[ "--help" == "${fileName}" ]]||[[ "-h" == "${fileName}" ]]; then
        log "gitresetfile FILENAME BRANCH"
        return
    fi
    if [[ -f "$fileName" ]]; then
        if [[ -z "$branch" ]]; then
            branch="origin/master"
        fi
        evaladvanced "git checkout $branch '$fileName'"
    else
        errorlog "Invalid file - $fileName"
    fi
}
function gitrepobackup {
	local url="$1"
	git clone --mirror "$url"
}
function gitreporestorebackup {
	local url="$1"
	git push --mirror "$url"
}
alias gitundolastcommit="evaladvanced 'git reset --soft HEAD~1'"
function gitmovsubmodule {
    local old="$1"
    local new="$2"
    local newParentDir="$(dirname "$new")"
    mkdir -p "$newParentDir"
    git mv "$old" "$new"
}
function gitaddscriptperm {
    local script="$1"
    local scriptFilename="$(basename "$script")"
    git update-index --chmod=+x "$script"
    git ls-files --stage | grep "$scriptFilename"
}
function gitcherrypickmaster {
    local commit="$1"
    git cherry-pick -m 1 "$commit"
}
function gitcherrypickmastercontinue {
    git cherry-pick --continue
}
alias gitclone="git clone"
function githubchangeurl() {
    read -p "Github Username: " username
    read -p "Github Token: " token
    read -p "Github URL end path(ex: AAA/bbb.gi): " urlEndPath
    local url="https://${username}:${token}@github.com/$urlEndPath"
    infolog "Set new github URL: $url"
    git remote set-url origin "$url"
}
function gitsetconfig() {
    local -a configArr=(
        "core.autocrlf input"
        "core.fileMode false"
        "core.logAllRefUpdates true"
        "core.ignorecase true"
        "pull.rebase true"
        "--unset safe.directory"
        "--add safe.directory '*'"
        "merge.ff false"
    )
    for configCmd in "${configArr[@]}"; do
        evaladvanced "git config --global $configCmd"
    done
    if [ $(directoryexists "$PWD/.git") = "true" ]||[ $(fileexists "$PWD/.git") = "true" ]; then
        infolog "Set local configurations"
        for configCmd in "${configArr[@]}"; do
            evaladvanced "git config $configCmd"
        done
    fi
}
function gitconfiguser() {
    read -p "Username: " username
    read -p "Email: " email
    local global_var="--global"
    if [[ -d "$PWD/.git" ]]||[[ -f "$PWD/.git" ]]; then
        global_var=""
    fi
    evaladvanced "git config ${global_var} user.name \"$username\""
    evaladvanced "git config ${global_var} user.email \"$email\""
}
alias gitcommit="git commit -m"
alias gitstageall="git add ."
alias gitstatus="git status"
function gitlatestversionrepo() {
    local owner="$1"
    local repo="$2"
    local isrelease=$3
    local urlsufix="/latest"
    if [[ -n "${isrelease}" ]]||[[ "${isrelease}" == "true" ]]; then
        urlsufix=""
    fi
    local url="https://api.github.com/repos/$owner/$repo/releases${urlsufix}"
    local version=$(curl -s "$url" | grep -Po '"tag_name": "\K[^"]*' | head -n 1)
    if [[ -z "${version}" ]]; then
        version=$(curl -s "$url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | head -n 1)
    fi
    if [[ "$version" == "v"* ]]; then
        version=$(echo "$version" | grep -Po 'v\K.*')
    fi
    echo "$version"
}
