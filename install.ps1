#Requires -Version 5.1
$ErrorActionPreference = 'Stop'

$Repo = "gurselcakar/arithmego"
$InstallDir = Join-Path $env:LOCALAPPDATA "arithmego"

function Fail($msg) {
    Write-Host "`n  $msg`n" -ForegroundColor Red
    exit 1
}

# Detect architecture
$arch = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture) {
    'X64'   { 'amd64' }
    'Arm64' { 'arm64' }
    default { Fail "Unsupported architecture: $_" }
}

Write-Host ""

# Fetch latest version
$spinChars = @('+', [char]0x2212, [char]0x00D7, [char]0x00F7)
$spinIdx = 0

Write-Host -NoNewline "  $($spinChars[$spinIdx])  Fetching latest version..."
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
    $version = $release.tag_name
} catch {
    Fail "Could not determine latest version"
}
if (-not $version) { Fail "Could not determine latest version" }
Write-Host "`r  +  Fetching latest version... done" -NoNewline
Write-Host ""

# Download
$filename = "arithmego_windows_${arch}.zip"
$url = "https://github.com/$Repo/releases/download/$version/$filename"
$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) "arithmego-install-$([guid]::NewGuid().ToString('N').Substring(0,8))"
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
$zipPath = Join-Path $tmpDir $filename

Write-Host -NoNewline "  $($spinChars[1])  Downloading arithmego $version..."
try {
    Invoke-WebRequest -Uri $url -OutFile $zipPath -UseBasicParsing
} catch {
    Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    Fail "Download failed - check if the release exists for windows/$arch"
}
Write-Host "`r  +  Downloading arithmego $version... done" -NoNewline
Write-Host ""

# Extract
try {
    Expand-Archive -Path $zipPath -DestinationPath $tmpDir -Force
} catch {
    Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    Fail "Failed to extract archive"
}

# Install
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
$exeSrc = Join-Path $tmpDir "arithmego.exe"
if (-not (Test-Path $exeSrc)) {
    Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    Fail "arithmego.exe not found in archive"
}
Copy-Item $exeSrc (Join-Path $InstallDir "arithmego.exe") -Force
Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "  ArithmeGo $version" -ForegroundColor Gray
Write-Host "  Installed to $InstallDir" -ForegroundColor Gray

# Add to PATH if needed
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($userPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$InstallDir", 'User')
    $env:Path = "$env:Path;$InstallDir"
    Write-Host ""
    Write-Host "  Added to PATH. Restart your terminal for it to take effect." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "  Run " -NoNewline
Write-Host "arithmego" -NoNewline -ForegroundColor White
Write-Host " to start playing!"
Write-Host ""
Write-Host "  --" -ForegroundColor DarkGray
Write-Host "  Your AI is thinking. You should too." -ForegroundColor DarkGray
Write-Host "  Built by Gursel Cakar." -ForegroundColor DarkGray
Write-Host ""
