# Author: José M. C. Noronha
[cmdletbinding()]
param (
    [string] $appId,
    [string] $title,
    [string]
    [parameter(ValueFromPipeline)]
    $message,
    [string] $icon,
    [Alias("h")]
    [switch] $help
)

if ($help) {
    Write-Host "$PSCommandPath -appId ID -title TITLE -message MESSAGE -icon ICON"
    exit 0
}

[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$Template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastImageAndText02)

$RawXml = [xml] $Template.GetXml()
($RawXml.toast.visual.binding.text|Where-Object {$_.id -eq "1"}).AppendChild($RawXml.CreateTextNode($title)) > $null
($RawXml.toast.visual.binding.text|Where-Object {$_.id -eq "2"}).AppendChild($RawXml.CreateTextNode($message)) > $null
if (![string]::IsNullOrEmpty($icon)) {
    ($RawXml.toast.visual.binding.image|Where-Object {$_.id -eq "1"}).SetAttribute('src', $icon) > $null
}

$SerializedXml = New-Object Windows.Data.Xml.Dom.XmlDocument
$SerializedXml.LoadXml($RawXml.OuterXml)
Write-Host $SerializedXml.ToString()

$Toast = [Windows.UI.Notifications.ToastNotification]::new($SerializedXml)
$Toast.Tag = "$appId"
$Toast.Group = "$appId"
$Toast.ExpirationTime = [DateTimeOffset]::Now.AddMinutes(1)

$Notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("$appId")
$Notifier.Show($Toast);