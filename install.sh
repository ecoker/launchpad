#!/bin/sh
# Launchpad installer
# Usage: curl -fsSL https://raw.githubusercontent.com/ehrencoker/agent-kit/main/install.sh | sh
#
# This script downloads the latest launchpad binary for your platform
# and installs it to /usr/local/bin (or a custom location).

set -e

REPO="ehrencoker/agent-kit"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY="launchpad"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux)  PLATFORM="linux" ;;
  darwin) PLATFORM="darwin" ;;
  *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Colors
CYAN="\033[36m"
GREEN="\033[32m"
MAGENTA="\033[35m"
DIM="\033[2m"
BOLD="\033[1m"
RESET="\033[0m"

echo ""
echo "${CYAN}${BOLD}   🚀 launchpad installer${RESET}"
echo ""

# Get latest release tag
echo "${DIM}Finding latest release...${RESET}"
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Could not determine latest version. Check https://github.com/${REPO}/releases"
  exit 1
fi

VERSION="${LATEST#v}"
FILENAME="${BINARY}_${PLATFORM}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

echo "${DIM}Downloading ${BINARY} ${VERSION} for ${PLATFORM}/${ARCH}...${RESET}"
TMPDIR=$(mktemp -d)
curl -fsSL "$URL" -o "${TMPDIR}/${FILENAME}"

echo "${DIM}Extracting...${RESET}"
tar -xzf "${TMPDIR}/${FILENAME}" -C "${TMPDIR}"

echo "${DIM}Installing to ${INSTALL_DIR}...${RESET}"
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi
chmod +x "${INSTALL_DIR}/${BINARY}"

rm -rf "$TMPDIR"

echo ""
echo "${GREEN}\u2714${RESET} ${BOLD}launchpad ${MAGENTA}v${VERSION}${RESET} installed to ${CYAN}${INSTALL_DIR}/${BINARY}${RESET}"
echo ""
echo "${DIM}Get started:${RESET}"
echo "  ${CYAN}launchpad init ./my-app${RESET}"
echo ""
