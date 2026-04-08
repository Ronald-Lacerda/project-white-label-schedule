$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$GoBin = "C:\Program Files\Go\bin\go.exe"

Push-Location (Join-Path $RootDir "backend")
try {
    & $GoBin run ./cmd/api
}
finally {
    Pop-Location
}
