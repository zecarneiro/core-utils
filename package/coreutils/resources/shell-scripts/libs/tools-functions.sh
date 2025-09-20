#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2164
# shellcheck disable=SC2155

# Feature for linux only
alias restart-pipewire="systemctl --user restart pipewire.service"
function changedefaultjdk {
    local java_default_script_name="/etc/profile.d/jdk-default.sh"
    evaladvanced "update-java-alternatives --list"
    read -p "Insert Path of java(JAVA_HOME): " javaHome
    evaladvanced "echo \"JAVA_HOME_DEFAULT='${javaHome}'\" | sudo tee ${java_default_script_name}"
    evaladvanced "echo \"export JAVA_HOME=${javaHome}\" | sudo tee -a ${java_default_script_name}"
    evaladvanced "echo \"export PATH=\\$PATH:\\${JAVA_HOME_DEFAULT}/bin\" | sudo tee -a ${java_default_script_name}"
    evaladvanced "source ${java_default_script_name}"
}
