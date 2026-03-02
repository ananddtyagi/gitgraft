#!/bin/bash
set -e

# Git-Graft Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/anandtyagi/gitgraft/main/install.sh | bash

REPO="anandtyagi/gitgraft"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="graft"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() { echo -e "${BLUE}==>${NC} $1"; }
success() { echo -e "${GREEN}==>${NC} $1"; }
error() { echo -e "${RED}Error:${NC} $1" >&2; exit 1; }

# Detect OS and architecture
detect_platform() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac

    case "$OS" in
        darwin|linux) ;;
        mingw*|msys*|cygwin*) OS="windows" ;;
        *) error "Unsupported OS: $OS" ;;
    esac

    PLATFORM="${OS}-${ARCH}"
}

# Check if Go is installed
has_go() {
    command -v go &> /dev/null
}

# Install via Go
install_with_go() {
    info "Installing with Go..."
    go install "github.com/${REPO}/cmd/graft@latest"

    # Find where it was installed
    GOBIN="${GOBIN:-$(go env GOPATH)/bin}"
    if [[ -f "${GOBIN}/graft" ]]; then
        success "Installed to ${GOBIN}/graft"
        echo ""
        echo "Make sure ${GOBIN} is in your PATH:"
        echo "  export PATH=\"\$PATH:${GOBIN}\""
    fi
}

# Install from pre-built binary (for future releases)
install_from_release() {
    info "Downloading pre-built binary..."

    # Get latest release
    LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

    if [[ -z "$LATEST" ]]; then
        error "Could not find latest release. Try installing with Go instead."
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/graft-${PLATFORM}"
    [[ "$OS" == "windows" ]] && DOWNLOAD_URL="${DOWNLOAD_URL}.exe"

    # Download
    TMP_FILE=$(mktemp)
    curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE" || error "Download failed"

    # Install
    chmod +x "$TMP_FILE"

    if [[ -w "$INSTALL_DIR" ]]; then
        mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        info "Need sudo to install to ${INSTALL_DIR}"
        sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    success "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Build from source
build_from_source() {
    info "Building from source..."

    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    git clone --depth 1 "https://github.com/${REPO}.git" .
    go build -o graft ./cmd/graft

    if [[ -w "$INSTALL_DIR" ]]; then
        mv graft "${INSTALL_DIR}/${BINARY_NAME}"
    else
        info "Need sudo to install to ${INSTALL_DIR}"
        sudo mv graft "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    cd - > /dev/null
    rm -rf "$TMP_DIR"

    success "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Setup shell alias
setup_alias() {
    echo ""
    read -p "Set up 'gg' alias for quick access? [Y/n] " -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Nn]$ ]]; then
        return
    fi

    SHELL_NAME=$(basename "$SHELL")
    case "$SHELL_NAME" in
        zsh)  RC_FILE="$HOME/.zshrc" ;;
        bash) RC_FILE="$HOME/.bashrc" ;;
        fish) RC_FILE="$HOME/.config/fish/config.fish" ;;
        *)    RC_FILE="" ;;
    esac

    if [[ -n "$RC_FILE" ]]; then
        echo "" >> "$RC_FILE"
        echo "# Git-Graft alias" >> "$RC_FILE"
        echo "alias gg='graft'" >> "$RC_FILE"
        success "Added 'gg' alias to $RC_FILE"
        echo "  Run: source $RC_FILE"
    else
        echo "Add this to your shell config:"
        echo "  alias gg='graft'"
    fi
}

# Main
main() {
    echo ""
    echo "  ╔═╗╦╔╦╗  ╔═╗╦═╗╔═╗╔═╗╔╦╗"
    echo "  ║ ╦║ ║───║ ╦╠╦╝╠═╣╠╣  ║ "
    echo "  ╚═╝╩ ╩   ╚═╝╩╚═╩ ╩╚   ╩ "
    echo ""
    echo "  A chill TUI for git"
    echo ""

    detect_platform
    info "Detected platform: $PLATFORM"

    if has_go; then
        install_with_go
    else
        info "Go not found, building from source requires Go"
        echo ""
        echo "Install Go first: https://go.dev/dl/"
        echo "Or install with Homebrew: brew install go"
        exit 1
    fi

    setup_alias

    echo ""
    success "Installation complete!"
    echo ""
    echo "  Run 'graft' (or 'gg') in any git repository"
    echo ""
}

main "$@"
