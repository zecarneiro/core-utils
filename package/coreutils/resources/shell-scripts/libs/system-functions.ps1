# Author: José M. C. Noronha
# Some code has source: https://github.com/ChrisTitusTech/powershell-profile

function setenv {
  param(
    [string] $envKey,
    [string] $envValue,
    [ValidateSet('Machine','User')]
    [string] $envType,
    [Alias("h")]
    [switch] $help
  )
  if ($help) {
    log "setenv ENV_KEY ENV_VALUE ENV_TYPE"
    return
  }
  if ([string]::IsNullOrEmpty($envKey) -or [string]::IsNullOrEmpty($envType)) {
    errorlog "Invalid envKey or envType"
  } else {
    if ($envType.Equals("Machine")) {
      $envType = [System.EnvironmentVariableTarget]::Machine
    } else {
      $envType = [System.EnvironmentVariableTarget]::User
    }
    [Environment]::SetEnvironmentVariable("$envKey", "$envValue", "$envType")
  }
}
function deleteenv {
  param(
    [string] $envKey,
    [ValidateSet('Machine','User')]
    [string] $envType,
    [Alias("h")]
    [switch] $help
  )
  if ($help) {
    log "deleteenv ENV_KEY ENV_TYPE"
    return
  }
  if ([string]::IsNullOrEmpty($envKey) -or [string]::IsNullOrEmpty($envType)) {
    errorlog "Invalid envKey or envType"
  } else {
    if ($envType.Equals("Machine")) {
      $envType = [System.EnvironmentVariableTarget]::Machine
    } else {
      $envType = [System.EnvironmentVariableTarget]::User
    }
    [Environment]::SetEnvironmentVariable("$envKey", [NullString]::Value, "$envType")
  }
}

# Features that already exist on linux
function export($expression) {
  if ([string]::IsNullOrEmpty($expression)) {
    Get-ChildItem env:*
  }
  else {
    if (!$expression.Contains("=")) {
      Write-Output "Environment variable $expression not defined"
    }
    else {
      $expressionArr = $expression.Split("=")
      $name = $expressionArr[0]
      $value = ""
      if ($expressionArr.Length -gt 1) {
        $value = $expressionArr[1]
      }
      set-item -force -path "env:$name" -value $value
    }
  }
}
function removeduplicatedenvval {
  param(
    [string] $envKey,
    [ValidateSet('Machine','User')]
    [string] $envType,
    [Alias("h")]
    [switch] $help
  )
  if ($help) {
    log "removeduplicatedenvval ENV_KEY ENV_TYPE"
    return
  }
  if ([string]::IsNullOrEmpty($envKey) -or [string]::IsNullOrEmpty($envType)) {
    errorlog "Invalid envKey or envType"
  } else {
    if ($envType.Equals("Machine")) {
      $envType = [System.EnvironmentVariableTarget]::Machine
    } else {
      $envType = [System.EnvironmentVariableTarget]::User
    }
    [Environment]::GetEnvironmentVariable("$envKey", $envType)
    $noDupesPath = (([Environment]::GetEnvironmentVariable("$envKey", $envType) -split ';' | Select-Object -Unique) -join ';')
    [Environment]::SetEnvironmentVariable("$envKey", $noDupesPath, $envType)
  }
}
