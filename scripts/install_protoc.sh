#!/bin/bash

set -e

# Define version and installation directory
PROTOC_VERSION="30.2"
INSTALL_DIR="/usr/local"

# Determine system architecture
ARCH=$(uname -m)
case $ARCH in
  x86_64) ARCH="x86_64" ;;
  aarch64) ARCH="aarch_64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Create a temporary directory for download
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download the specified version
ZIP_FILE="protoc-${PROTOC_VERSION}-linux-${ARCH}.zip"
DOWNLOAD_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${ZIP_FILE}"

echo "Downloading protoc ${PROTOC_VERSION} from ${DOWNLOAD_URL}..."
curl -LO "$DOWNLOAD_URL"

# Extract the zip file
unzip "$ZIP_FILE" -d protoc

# Install protoc binary and include files
echo "Installing protoc to ${INSTALL_DIR}..."
sudo mv protoc/bin/protoc "${INSTALL_DIR}/bin/"
sudo chmod +x "${INSTALL_DIR}/bin/protoc"
sudo cp -r protoc/include/* "${INSTALL_DIR}/include/"

# Clean up
cd ~
rm -rf "$TMP_DIR"

# Verify installation
echo "protoc installation completed."
protoc --version
