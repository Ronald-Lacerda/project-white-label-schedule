$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BackendScript = Join-Path $RootDir "run-backend.ps1"
$FrontendScript = Join-Path $RootDir "run-frontend.ps1"

$backendJob = Start-Job -ScriptBlock {
    param($ScriptPath)
    & powershell -NoProfile -ExecutionPolicy Bypass -File $ScriptPath
} -ArgumentList $BackendScript

try {
    & powershell -NoProfile -ExecutionPolicy Bypass -File $FrontendScript
}
finally {
    if ($backendJob.State -eq "Running") {
        Stop-Job -Job $backendJob
    }

    Receive-Job -Job $backendJob -Keep | Out-Null
    Remove-Job -Job $backendJob -Force
}
