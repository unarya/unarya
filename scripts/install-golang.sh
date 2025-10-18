#!/usr/bin/env bash
set -e

GO_VERSION=1.25.0
GO_TAR=go${GO_VERSION}.linux-amd64.tar.gz
INSTALL_DIR=/usr/local
PROFILE_FILE="$HOME/.bashrc"

echo "=== Installing Go ${GO_VERSION} ==="

# Remove old Go if exists
if [ -d "${INSTALL_DIR}/go" ]; then
  echo "Removing old Go installation..."
  sudo rm -rf ${INSTALL_DIR}/go
fi

# Download new version
echo "Downloading Go ${GO_VERSION} ..."
wget -q https://go.dev/dl/${GO_TAR} -O /tmp/${GO_TAR}

# Extract to /usr/local
echo "Extracting to ${INSTALL_DIR} ..."
sudo tar -C ${INSTALL_DIR} -xzf /tmp/${GO_TAR}
rm /tmp/${GO_TAR}

# Add Go to PATH if missing
if ! grep -q "/usr/local/go/bin" "$PROFILE_FILE"; then
  echo "Adding Go to PATH in $PROFILE_FILE ..."
  echo 'export PATH=$PATH:/usr/local/go/bin' >> "$PROFILE_FILE"
fi

# Apply PATH immediately
export PATH=$PATH:/usr/local/go/bin

# Verify installation
echo "Verifying installation ..."
go version || { echo "❌ Go installation failed."; exit 1; }

echo "✅ Go ${GO_VERSION} installed successfully!"

