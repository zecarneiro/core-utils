# Author: José M. C. Noronha

# Load Assembly
Add-Type -AssemblyName System.Windows.Forms

$OpenFileDialog = New-Object System.Windows.Forms.OpenFileDialog
$OpenFileDialog.ShowDialog() | Out-Null
$filename = $OpenFileDialog.FileName
return @{ selected="$filename"; }
