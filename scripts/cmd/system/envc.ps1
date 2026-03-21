# Author: José M. C. Noronha
# Some code has source: https://github.com/ChrisTitusTech/powershell-profile

param(
  [Parameter(Mandatory=$true)]
  [string] $envKey,
  [Parameter(Mandatory=$true)]
  [ValidateSet('Machine','User')]
  [string] $envType,
  [Parameter(Mandatory=$true)]
  [ValidateSet('get','set','delete','clean')]
  [string] $operation,
  [string[]] $envValue,
  [Alias("h")]
  [switch] $help
)

function GetEnv {
  if ($help) {
    log "GetEnv ENV_KEY ENV_TYPE"
    return
  }
  if ([string]::IsNullOrEmpty($envKey) -or [string]::IsNullOrEmpty($envType)) {
    error-log "Invalid envKey or envType"
  } else {
    if ($envType.Equals("Machine")) {
      $envType = [System.EnvironmentVariableTarget]::Machine
    } else {
      $envType = [System.EnvironmentVariableTarget]::User
    }
    ([Environment]::GetEnvironmentVariable("$envKey", $envType)) | Out-String
  }
}

# This function will replace old env if already exists
function SetEnv {
  if ($help) {
    log "SetEnv ENV_KEY ENV_VALUE ENV_TYPE"
    return
  }
  if ([string]::IsNullOrEmpty($envKey) -or [string]::IsNullOrEmpty($envType) -or $null -eq $envValue) {
    error-log "Invalid envKey or envType or envValue"
  } else {
    if ($envType.Equals("Machine")) {
      $envType = [System.EnvironmentVariableTarget]::Machine
    } else {
      $envType = [System.EnvironmentVariableTarget]::User
    }
    if ($envValue.Count -eq 1 -and $envValue -match ',') {
      $envValue = $envValue -split ','
    }
    $envValueString = $envValue -join ';'
    [Environment]::SetEnvironmentVariable("$envKey", "$envValueString", "$envType")
  }
}

# This function will delete env from the system
function DeleteEnv {
  if ($help) {
    log "DeleteEnv ENV_KEY ENV_TYPE"
    return
  }
  if ([string]::IsNullOrEmpty($envKey) -or [string]::IsNullOrEmpty($envType)) {
    error-log "Invalid envKey or envType"
  } else {
    if ($envType.Equals("Machine")) {
      $envType = [System.EnvironmentVariableTarget]::Machine
    } else {
      $envType = [System.EnvironmentVariableTarget]::User
    }
    [Environment]::SetEnvironmentVariable("$envKey", [NullString]::Value, "$envType")
  }
}

function CleanEnv {
  if ($help) {
    log "CleanEnv ENV_KEY ENV_TYPE"
    return
  }
  if ([string]::IsNullOrEmpty($envKey) -or [string]::IsNullOrEmpty($envType)) {
    error-log "Invalid envKey or envType"
  } else {
    if ($envType.Equals("Machine")) {
      $envType = [System.EnvironmentVariableTarget]::Machine
    } else {
      $envType = [System.EnvironmentVariableTarget]::User
    }
    $envKeyValue = ([Environment]::GetEnvironmentVariable("$envKey", $envType))
    $noDupesPath = (($envKeyValue -split ';' | Select-Object -Unique) -join ';')
    [Environment]::SetEnvironmentVariable("$envKey", $noDupesPath, $envType)
  }
}

switch ($operation) {
    "get"  { GetEnv }
    "set" { SetEnv }
    "delete" { DeleteEnv }
    "clean" { CleanEnv }
    Default  { "Need to pass a operation" }
}
