$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$FrontendPort = 3000

function Stop-ListenersOnPort($Port) {
    $pids = @()

    try {
        $listeners = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction Stop
        $pids = $listeners | Select-Object -ExpandProperty OwningProcess -Unique
    }
    catch {
        $netstatLines = netstat -ano | Select-String ":$Port"
        foreach ($line in $netstatLines) {
            $parts = ($line.ToString() -split '\s+') | Where-Object { $_ }
            if ($parts.Length -ge 5 -and $parts[3] -eq 'LISTENING') {
                $listenerPid = $parts[-1]
                if ($listenerPid -match '^\d+$') {
                    $pids += [int]$listenerPid
                }
            }
        }
        $pids = $pids | Select-Object -Unique
    }

    foreach ($listenerPid in $pids) {
        try {
            Stop-Process -Id $listenerPid -Force -ErrorAction Stop
        }
        catch {
            Write-Warning "Nao foi possivel encerrar o processo $listenerPid na porta $Port."
        }
    }
}

Stop-ListenersOnPort $FrontendPort

Push-Location (Join-Path $RootDir "frontend")
try {
    npm run dev
}
finally {
    Pop-Location
}
