# PowerShell setup script for Eventus devs
Set-Location -Path (Split-Path -Parent $MyInvocation.MyCommand.Path)
Write-Host "Running go mod tidy..."
go mod tidy
Write-Host "Running unit tests..."
go test ./... -v
Write-Host "Done."
