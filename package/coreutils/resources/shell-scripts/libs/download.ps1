param (
    # Parameter help description
    [Parameter(Mandatory)]
    [string]$url,

    # Parameter help description
    [Parameter(Mandatory)]
    [Alias("o")]
    [string]$file
)
Begin {
    function Show-Progress {
        param (
            # Enter total value
            [Parameter(Mandatory)]
            [Single]$TotalValue,

            # Enter current value
            [Parameter(Mandatory)]
            [Single]$CurrentValue,

            # Enter custom progresstext
            [Parameter(Mandatory)]
            [string]$ProgressText,

            # Enter value suffix
            [Parameter()]
            [string]$ValueSuffix,

            # Enter bar lengh suffix
            [Parameter()]
            [int]$BarSize = 40,

            # show complete bar
            [Parameter()]
            [switch]$Complete
        )

        # calc %
        $percent = $CurrentValue / $TotalValue
        $percentComplete = $percent * 100
        if ($ValueSuffix) {
            $ValueSuffix = " $ValueSuffix" # add space in front
        }
        if ($psISE) {
            Write-Progress "$ProgressText $CurrentValue$ValueSuffix of $TotalValue$ValueSuffix" -id 0 -percentComplete $percentComplete
        }
        else {
            # build progressbar with string function
            $curBarSize = $BarSize * $percent
            $progbar = ""
            $progbar = $progbar.PadRight($curBarSize, [char]9608)
            $progbar = $progbar.PadRight($BarSize, [char]9617)

            if (!$Complete.IsPresent) {
                Write-Host -NoNewLine "`r$ProgressText $progbar [ $($CurrentValue.ToString("#.###").PadLeft($TotalValue.ToString("#.###").Length))$ValueSuffix / $($TotalValue.ToString("#.###"))$ValueSuffix ] $($percentComplete.ToString("##0.00").PadLeft(6)) % complete"
            }
            else {
                Write-Host -NoNewLine "`r$ProgressText $progbar [ $($TotalValue.ToString("#.###").PadLeft($TotalValue.ToString("#.###").Length))$ValueSuffix / $($TotalValue.ToString("#.###"))$ValueSuffix ] $($percentComplete.ToString("##0.00").PadLeft(6)) % complete"
            }
        }
    }
}
Process {
    try {
        $storeEAP = $ErrorActionPreference
        $ErrorActionPreference = 'Stop'

        # invoke request
        $request = [System.Net.HttpWebRequest]::Create($url)
        $response = $request.GetResponse()

        if ($response.StatusCode -eq 401 -or $response.StatusCode -eq 403 -or $response.StatusCode -eq 404) {
            throw "Remote file either doesn't exist, is unauthorized, or is forbidden for '$url'."
        }

        if ($file -match '^\.\\') {
            $file = Join-Path (Get-Location -PSProvider "FileSystem") ($file -Split '^\.')[1]
        }

        if ($file -and !(Split-Path $file)) {
            $file = Join-Path (Get-Location -PSProvider "FileSystem") $file
        }

        if ($file) {
            $fileDirectory = $([System.IO.Path]::GetDirectoryName($file))
            if (!(Test-Path($fileDirectory))) {
                [System.IO.Directory]::CreateDirectory($fileDirectory) | Out-Null
            }
        }

        [long]$fullSize = $response.ContentLength
        $fullSizeMB = $fullSize / 1024 / 1024

        # define buffer
        [byte[]]$buffer = new-object byte[] 1048576
        [long]$total = [long]$count = 0

        # create reader / writer
        $reader = $response.GetResponseStream()
        $writer = new-object System.IO.FileStream $file, "Create"
        $fileProgress = (filename "$file")

        # start download
        $finalBarCount = 0 #show final bar only one time
        do {

            $count = $reader.Read($buffer, 0, $buffer.Length)

            $writer.Write($buffer, 0, $count)

            $total += $count
            $totalMB = $total / 1024 / 1024

            if ($fullSize -gt 0) {
                Show-Progress -TotalValue $fullSizeMB -CurrentValue $totalMB -ProgressText "$($fileProgress)" -ValueSuffix "MB"
            }

            if ($total -eq $fullSize -and $count -eq 0 -and $finalBarCount -eq 0) {
                Show-Progress -TotalValue $fullSizeMB -CurrentValue $totalMB -ProgressText "$($fileProgress)" -ValueSuffix "MB" -Complete
                $finalBarCount++
            }
        } while ($count -gt 0)
        Write-Host -NoNewLine "`n"
    }
    catch {
        $ExeptionMsg = $_.Exception.Message
        error-log "$ExeptionMsg"
    }

    finally {
        # cleanup
        if ($reader) { $reader.Close() }
        if ($writer) { $writer.Flush(); $writer.Close() }

        $ErrorActionPreference = $storeEAP
        [GC]::Collect()
    }
}