#!/bin/sh
set -e

REPO="nopecho/claude-television"
BIN_NAME="ctv"

# OS detection
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS_NAME="Linux" ;;
    Darwin*)    OS_NAME="Darwin" ;;
    *)          echo "Unsupported OS: ${OS}"; exit 1 ;;
esac

# Arch detection
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)  ARCH_NAME="x86_64" ;;
    amd64)   ARCH_NAME="x86_64" ;;
    arm64)   ARCH_NAME="arm64" ;;
    aarch64) ARCH_NAME="arm64" ;;
    *)       echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
esac

echo "Detecting latest release for ${REPO}..."
# Get the redirect URL of the latest release
LATEST_URL=$(curl -w "%{url_effective}\n" -I -L -s -S "https://github.com/${REPO}/releases/latest" -o /dev/null)
VERSION=$(basename "${LATEST_URL}")

if [ -z "${VERSION}" ] || [ "${VERSION}" = "latest" ]; then
    echo "Failed to determine latest version."
    exit 1
fi

echo "Latest version is ${VERSION}"

# Construct the filename based on goreleaser template
TAR_FILE="${BIN_NAME}_${OS_NAME}_${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${TAR_FILE}"

TMP_DIR=$(mktemp -d)
# Ensure cleanup on exit
trap 'rm -rf "${TMP_DIR}"' EXIT

echo "Downloading ${TAR_FILE}..."
if ! curl -fL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${TAR_FILE}"; then
    echo "Failed to download ${DOWNLOAD_URL}"
    exit 1
fi

echo "Extracting..."
tar -xzf "${TMP_DIR}/${TAR_FILE}" -C "${TMP_DIR}"

INSTALL_DIR="/usr/local/bin"

echo "Installing to ${INSTALL_DIR}..."
if [ ! -w "${INSTALL_DIR}" ]; then
    echo "Root privileges are required to install to ${INSTALL_DIR}."
    sudo mv "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
    sudo chmod +x "${INSTALL_DIR}/${BIN_NAME}"
else
    mv "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
    chmod +x "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "✅ Successfully installed ${BIN_NAME} to ${INSTALL_DIR}/${BIN_NAME}"
echo "Run 'ctv' to get started."
