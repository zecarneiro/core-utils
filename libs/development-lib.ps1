function install-development-package {
    __install_node_typescript_javascript
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

function __install_node_typescript_javascript {
    if ((__show_install_message_question "NodeJS/Javascript/Typescript") -eq "y") {
        evaladvanced "scoop bucket add main"
        evaladvanced "scoop install main/nodejs-lts"
        . reloadprofile
        evaladvanced "npm install -g typescript"
    }
}

function __install_python {
    if ((__show_install_message_question "Python3/PIP") -eq "y") {
        evaladvanced "scoop bucket add main"
        evaladvanced "scoop install main/python"
        evaladvanced "pip install virtualenv"
    }
}

function __install_java {
    $message = "`nSet JAVA_HOME in option."
    # To download executable go to: https://adoptopenjdk.net/ or https://adoptium.net/
    if ((__show_install_message_question "Java JDK 17") -eq "y") {
        infolog "$message"
        evaladvanced "winget install -i --id=EclipseAdoptium.Temurin.17.JDK --accept-source-agreements --accept-package-agreements"
    }
    if ((__show_install_message_question "Java JDK 21") -eq "y") {
        infolog "$message"
        evaladvanced "winget install -i --id=EclipseAdoptium.Temurin.21.JDK --accept-source-agreements --accept-package-agreements"
    }
}

function __install_maven {
    if ((__show_install_message_question "Maven") -eq "y") {
        evaladvanced "scoop bucket add main"
        evaladvanced "scoop install main/maven"
    }
}

function __install_cpp_c {
    if ((__show_install_message_question "C/C++/Make/CLang") -eq "y") {
        log "`nAdd PATH for LLVM and CMake"
        evaladvanced "scoop bucket add main"
        evaladvanced "scoop install main/make"
        evaladvanced "scoop install main/gcc"
        evaladvanced "scoop install main/cmake"
        evaladvanced "scoop install main/clangd"
    }
}

function __install_php {
    if ((__show_install_message_question "PHP") -eq "y") {
        evaladvanced "scoop bucket add main"
        evaladvanced "scoop install main/php"
    }
}

function __install_golang {
    if ((__show_install_message_question "Go") -eq "y") {
        evaladvanced "scoop bucket add main"
        evaladvanced "scoop install main/go"
        . reloadprofile
        evaladvanced "go install golang.org/x/tools/gopls@latest"
        addalias "goclean" -command "go clean -cache -modcache -testcache -fuzzcache"
    }
}

function __install_sqlite3 {
    if ((__show_install_message_question "Sqlite3") -eq "y") {
        infolog "`nDownload link example: https://www.sqlite.org/2022/sqlite-tools-win32-x86-{version}.zip"
        evaladvanced "winget install --id=SQLite.SQLite --accept-source-agreements --accept-package-agreements"
    }
}

function __install_postgres_sql {
    if ((__show_install_message_question "Postgres SQL") -eq "y") {
        infolog "`nDownload link example: https://www.sqlite.org/2022/sqlite-tools-win32-x86-{version}.zip"
        infolog "For Client only only keep the options: Command line Tools"
        evaladvanced "winget install -i --id=PostgreSQL.PostgreSQL --accept-source-agreements --accept-package-agreements"
        & "$IMAGE_UTILS_DIR\postgressql.png"
    }
}

function __install_shell_language_server {
    if ((__show_install_message_question "Shellcheck") -eq "y") {
        evaladvanced "scoop install shellcheck"
    }
    if ((__show_install_message_question "Bash language server") -eq "y") {
        evaladvanced "npm install -g bash-language-server"
    }
}
