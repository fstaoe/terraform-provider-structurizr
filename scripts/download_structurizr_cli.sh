#!/bin/bash

set -e

# Structurizr CLI version to be downloaded
VERSION=$1
DIR=$2

# Checking if Structurizr CLI already exists
if [ -f "$DIR/structurizr.sh" ] && [ -f "$DIR/structurizr.bat" ]; then
    echo "Structurizr CLI already exists, skipping download."
    exit 0
fi

# Download the Structurizr CLI
curl -sL -o structurizr-cli.zip "https://github.com/structurizr/cli/releases/download/$VERSION/structurizr-cli.zip"

# Creating the necessary directories if do not exist
if [ ! -d "$DIR" ]; then
  mkdir -p "$DIR";
fi

# Unzip the Structurizr CLI
unzip structurizr-cli.zip -d "$DIR"
rm structurizr-cli.zip

# Make sure the shell script is executable
chmod +x "$DIR/structurizr.sh"

# Make sure the batch file is executable on Windows
if [ -f "$DIR/structurizr.bat" ]; then
    chmod +x "$DIR/structurizr.bat"
fi

# Setting current version
echo "$VERSION" > "$DIR/version.txt"