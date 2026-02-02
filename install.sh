#!/bin/bash
set -e

REPO="gurselcakar/arithmego"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}==>${NC} $1"
}

warn() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

error() {
    echo -e "${RED}Error:${NC} $1"
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *)       error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64)  echo "amd64" ;;
        amd64)   echo "amd64" ;;
        arm64)   echo "arm64" ;;
        aarch64) echo "arm64" ;;
        *)       error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest release version
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
        grep '"tag_name":' |
        sed -E 's/.*"([^"]+)".*/\1/'
}

main() {
    info "Installing arithmego..."

    # Check for required tools
    command -v curl >/dev/null 2>&1 || error "curl is required but not installed"
    command -v tar >/dev/null 2>&1 || error "tar is required but not installed"

    OS=$(detect_os)
    ARCH=$(detect_arch)

    info "Detected: ${OS}/${ARCH}"

    # Get latest version
    info "Fetching latest release..."
    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        error "Could not determine latest version"
    fi

    info "Latest version: ${VERSION}"

    # Download
    FILENAME="arithmego_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

    info "Downloading ${FILENAME}..."

    TMPDIR=$(mktemp -d)
    trap 'rm -rf "${TMPDIR}"' EXIT

    curl -fsSL "${URL}" -o "${TMPDIR}/${FILENAME}" || error "Download failed. Check if the release exists for your platform."

    # Extract
    info "Extracting..."
    tar -xzf "${TMPDIR}/${FILENAME}" -C "${TMPDIR}"

    # Install
    info "Installing to ${INSTALL_DIR}..."

    if [ -w "${INSTALL_DIR}" ]; then
        mv "${TMPDIR}/arithmego" "${INSTALL_DIR}/arithmego"
    else
        sudo mv "${TMPDIR}/arithmego" "${INSTALL_DIR}/arithmego"
    fi

    chmod +x "${INSTALL_DIR}/arithmego"

    info "Successfully installed arithmego ${VERSION}"
    echo ""
    echo "Run 'arithmego' to start playing!"
}

main "$@"
