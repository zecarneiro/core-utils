param([string] $file)

shell = new-object -comobject "Shell.Application"
if ((file-exists "$file")) {
    $file = (Resolve-Path -LiteralPath "$file")
    $shell.Namespace(0).ParseName("$file").InvokeVerb("delete")
} elseif ((directory-exists "$file")) {
    $file = (Resolve-Path -Path "$file")
    $shell.Namespace(0).ParseName("$file").InvokeVerb("delete")
}