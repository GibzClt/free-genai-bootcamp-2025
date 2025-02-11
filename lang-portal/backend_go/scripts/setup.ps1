# Check if Go is installed
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Go..."
    # Download Go installer
    Invoke-WebRequest -Uri "https://go.dev/dl/go1.21.0.windows-amd64.msi" -OutFile "go1.21.0.windows-amd64.msi"
    # Install Go
    Start-Process -Wait -FilePath "msiexec.exe" -ArgumentList "/i go1.21.0.windows-amd64.msi /quiet"
    # Clean up
    Remove-Item "go1.21.0.windows-amd64.msi"
    # Add Go to PATH
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine")
}

# Check if SQLite is installed
if (!(Test-Path "C:\sqlite3\sqlite3.exe")) {
    Write-Host "Installing SQLite..."
    # Create SQLite directory
    New-Item -ItemType Directory -Force -Path "C:\sqlite3"
    # Download SQLite
    Invoke-WebRequest -Uri "https://www.sqlite.org/2024/sqlite-tools-win32-x86-3440200.zip" -OutFile "sqlite.zip"
    # Extract SQLite
    Expand-Archive -Path "sqlite.zip" -DestinationPath "C:\sqlite3" -Force
    # Clean up
    Remove-Item "sqlite.zip"
    # Add SQLite to PATH
    [Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\sqlite3", "Machine")
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine")
}

# Install Mage if not installed
if (!(Get-Command mage -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Mage..."
    go install github.com/magefile/mage@latest
}

# Download Go dependencies
Write-Host "Downloading Go dependencies..."
go mod download
go mod tidy

# Create necessary directories
New-Item -ItemType Directory -Force -Path "db/migrations"
New-Item -ItemType Directory -Force -Path "db/seeds"

# Initialize database and run migrations
Write-Host "Initializing database..."
mage initDB
mage migrate
mage seed

Write-Host "Setup complete!" 