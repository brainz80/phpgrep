param(
    [Parameter(Position = 0)]
    [ValidateSet('lint', 'test')]
    [string]$Target = 'lint',

    [string]$Image = 'golang:latest'
)

$ErrorActionPreference = 'Stop'

$root = (Get-Location).Path -replace '\\', '/'
if ($root -match '^[A-Za-z]:') {
    $drive = $root.Substring(0, 1).ToLower()
    $root = "/$drive$($root.Substring(2))"
}

docker run --rm -v "${root}:/app" -w /app $Image make $Target
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
