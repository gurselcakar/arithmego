#!/bin/bash
set -e

REPO="gurselcakar/arithmego"
INSTALL_DIR="$HOME/.local/bin"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
DIM='\033[2m'
BOLD='\033[1m'
NC='\033[0m'

cleanup() {
    tput cnorm 2>/dev/null || true
    [ -n "${TMPDIR:-}" ] && rm -rf "${TMPDIR}"
}
trap cleanup EXIT

fail() {
    printf "\n  ${RED}%s${NC}\n\n" "$1"
    exit 1
}

# Math-themed spinner: cycles through arithmetic operators
spin() {
    local pid=$1
    local msg=$2
    local frames=('+' '−' '×' '÷')
    local i=0
    tput civis 2>/dev/null || true
    while kill -0 "$pid" 2>/dev/null; do
        printf "\r  ${DIM}%s${NC}  %s" "${frames[$i]}" "$msg"
        i=$(( (i + 1) % ${#frames[@]} ))
        sleep 0.12
    done
    wait "$pid" 2>/dev/null && local ok=1 || local ok=0
    tput cnorm 2>/dev/null || true
    if [ "$ok" = "1" ]; then
        printf "\r\033[2K"
    else
        printf "\r  ${RED}%s${NC}   \n" "$msg"
        return 1
    fi
}

detect_platform() {
    case "$(uname -s)" in
        Linux*)  OS="linux" ;;
        Darwin*) OS="darwin" ;;
        *)       fail "Unsupported OS: $(uname -s)" ;;
    esac
    case "$(uname -m)" in
        x86_64|amd64)   ARCH="amd64" ;;
        arm64|aarch64)  ARCH="arm64" ;;
        *)              fail "Unsupported architecture: $(uname -m)" ;;
    esac
}

main() {
    command -v curl >/dev/null 2>&1 || fail "curl is required"
    command -v tar >/dev/null 2>&1 || fail "tar is required"

    echo ""

    detect_platform

    # Fetch latest version
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    [ -z "$VERSION" ] && fail "Could not determine latest version"

    FILENAME="arithmego_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"
    TMPDIR=$(mktemp -d)

    # Download and extract
    ( curl -fsSL "${URL}" -o "${TMPDIR}/${FILENAME}" && tar -xzf "${TMPDIR}/${FILENAME}" -C "${TMPDIR}" ) &
    spin $! "Downloading arithmego ${VERSION}" \
        || fail "Download failed — check if the release exists for ${OS}/${ARCH}"

    # Install
    mkdir -p "${INSTALL_DIR}"
    mv "${TMPDIR}/arithmego" "${INSTALL_DIR}/arithmego"
    chmod +x "${INSTALL_DIR}/arithmego"
    printf "  ArithmeGo ${DIM}%s${NC}\n" "${VERSION}"
    printf "  Installed to ${DIM}%s${NC}\n" "${INSTALL_DIR}"

    # PATH check
    if ! echo ":${PATH}:" | grep -q ":${INSTALL_DIR}:"; then
        echo ""
        printf "  ${YELLOW}Add ~/.local/bin to your PATH:${NC}\n"
        echo ""
        printf "     ${DIM}# bash${NC}\n"
        printf "     echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc\n"
        echo ""
        printf "     ${DIM}# zsh${NC}\n"
        printf "     echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc\n"
    fi

    echo ""
    printf "  Run ${BOLD}arithmego${NC} to start playing!\n"
    echo ""
    printf "  ${DIM}--${NC}\n"
    printf "  ${DIM}Your AI is thinking. You should too.${NC}\n"
    printf "  ${DIM}Built by Gürsel Çakar.${NC}\n"
    echo ""
}

main "$@"
