# core-utils

## Windows

1. Open terminal on this repository folder.
2. Run `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`
3. Run `$installStep=1; Invoke-RestMethod -Uri https://raw.githubusercontent.com/zecarneiro/core-utils/master/make.ps1 | Invoke-Expression`
4. **Restart Terminal**
4. Run `$installStep=2; Invoke-RestMethod -Uri https://raw.githubusercontent.com/zecarneiro/core-utils/master/make.ps1 | Invoke-Expression`

If you had any problem during installation, run this command on powershell:

`[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12`

## Linux - Bash