#!/bin/bash

# Install Go if not installed (for Linux)
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
fi

# Install SQLite if not installed
if ! command -v sqlite3 &> /dev/null; then
    echo "Installing SQLite..."
    sudo apt-get update
    sudo apt-get install -y sqlite3
fi

# Install Mage if not installed
if ! command -v mage &> /dev/null; then
    echo "Installing Mage..."
    go install github.com/magefile/mage@latest
fi

# Download Go dependencies
echo "Downloading Go dependencies..."
go mod download
go mod tidy

# Create necessary directories
mkdir -p db/migrations
mkdir -p db/seeds

# Initialize database and run migrations
echo "Initializing database..."
mage initDB
mage migrate
mage seed

echo "Setup complete!" 