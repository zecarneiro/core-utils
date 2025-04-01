#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2155

function install-development-package {
    __install_typescript
    __install_python
    __install_java
    __install_maven
    __install_cpp_c
    __install_php
    __install_golang
    __install_sqlite3
    __install_postgres_sql
    __install_shell_language_server
}

function __install_typescript {
    if [[ $(__show_install_message_question "Typescript") == "y" ]]; then
        install-npm
        # shellcheck source=/dev/null
        source "$HOME/.nvm/nvm.sh"
        evaladvanced "npm install -g typescript"
    fi
}

function __install_python {
    if [[ $(__show_install_message_question "Python3") == "y" ]]; then
        evaladvanced "sudo apt install python3 -y"
        evaladvanced "sudo apt install python-is-python3 -y"        
    fi
}

function __install_java {
    if [[ $(__show_install_message_question "Java JDK 17") == "y" ]]; then
        evaladvanced "sudo apt install openjdk-17-jdk -y"
        headerlog "When asked for path of java"
        headerlog "Insert: /usr/lib/jvm/java-17-openjdk-amd64"
        changedefaultjdk
    fi
    if [[ $(__show_install_message_question "Java JDK 21") == "y" ]]; then
        evaladvanced "sudo apt install openjdk-21-jdk -y"
        headerlog "When asked for path of java"
        headerlog "Insert: /usr/lib/jvm/java-21-openjdk-amd64"
        changedefaultjdk
    fi
}

function __install_maven {
    if [[ $(__show_install_message_question "Maven") == "y" ]]; then
        evaladvanced "sudo apt install maven -y"
    fi
}

function __install_cpp_c {
    if [[ $(__show_install_message_question "C/C++/Make/CLang/Objective-C") == "y" ]]; then
        evaladvanced "sudo apt install build-essential g++ gcc gdb cmake make clang clangd clang-format clang-tidy clang-tools -y"
    fi
}

function __install_php {
    if [[ $(__show_install_message_question "PHP") == "y" ]]; then
        evaladvanced "sudo apt install php -y"
    fi
}

function __install_golang {
    if [[ $(__show_install_message_question "Go") == "y" ]]; then
        evaladvanced "sudo apt install golang-go -y"
        reloadprofile
        evaladvanced "go install golang.org/x/tools/gopls@latest"
        addalias "goclean" "go clean -cache -modcache -testcache -fuzzcache"
        local bashrcFile="$HOME/.bashrc"
        if [ "$(filecontain "$bashrcFile" "/go/bin:")" == false ]; then
            writefile "$bashrcFile" "export PATH=\"\$HOME/go/bin:\$PATH\"" -append
        fi
    fi
}

function __install_sqlite3 {
    if [[ $(__show_install_message_question "Sqlite3") == "y" ]]; then
        infolog "\nDownload link example: https://www.sqlite.org/2022/sqlite-autoconf-{version}.tar.gz"
        evaladvanced "sudo apt install build-essential libsqlite3-dev -y"
        evaladvanced "sudo apt install sqlite3 -y"
    fi
}

function __install_postgres_sql {
    if [[ $(__show_install_message_question "Postgres SQL") == "y" ]]; then
        evaladvanced "sudo apt install curl ca-certificates -y"
        evaladvanced "sudo install -d /usr/share/postgresql-common/pgdg"
        evaladvanced "sudo curl -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc --fail https://www.postgresql.org/media/keys/ACCC4CF8.asc"
        evaladvanced "echo \"deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc] https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main\" | sudo tee /etc/apt/sources.list.d/pgdg.list >/dev/null"
        evaladvanced "sudo apt update"
        if [[ $(__show_install_message_question "Postgres SQL - Server") == "y" ]]; then
            evaladvanced "sudo apt install postgresql -y"
        fi
        if [[ $(__show_install_message_question "Postgres SQL - Client") == "y" ]]; then
            evaladvanced "sudo apt install postgresql-client -y"
        fi
    fi
}

function __install_shell_language_server {
    if [[ $(__show_install_message_question "Shellcheck") == "y" ]]; then
        evaladvanced "sudo apt install shellcheck -y"
    fi
    if [[ $(__show_install_message_question "Bash language server") == "y" ]]; then
        evaladvanced "npm install -g bash-language-server"
    fi
}
