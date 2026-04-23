#!/usr/bin/env bash
# bootstrap.sh — install all toolchains required for this monorepo.
# Idempotent: skips tools already present at the required version.
# Supports macOS (via Homebrew) and Debian/Ubuntu Linux.
set -euo pipefail

# ---------------------------------------------------------------------------
# helpers
# ---------------------------------------------------------------------------
info()  { printf '\033[0;34m[bootstrap] %s\033[0m\n' "$*"; }
ok()    { printf '\033[0;32m[bootstrap] %s\033[0m\n' "$*"; }
warn()  { printf '\033[0;33m[bootstrap] %s\033[0m\n' "$*"; }
die()   { printf '\033[0;31m[bootstrap] ERROR: %s\033[0m\n' "$*" >&2; exit 1; }

OS="$(uname -s)"
ARCH="$(uname -m)"

need_cmd() {
  command -v "$1" &>/dev/null
}

require_brew() {
  if ! need_cmd brew; then
    info "installing homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
}

# ---------------------------------------------------------------------------
# mise (manages node, go, rust versions via .tool-versions)
# ---------------------------------------------------------------------------
install_mise() {
  if need_cmd mise; then
    ok "mise already installed: $(mise --version)"
    return
  fi
  info "installing mise..."
  curl -fsSL https://mise.run | sh
  export PATH="$HOME/.local/bin:$PATH"
  ok "mise installed"
}

# ---------------------------------------------------------------------------
# Node.js via mise
# ---------------------------------------------------------------------------
install_node() {
  if need_cmd node; then
    ok "node already installed: $(node --version)"
  else
    info "installing node 22 via mise..."
    mise install node@22
    mise use -g node@22
    ok "node installed"
  fi
}

# ---------------------------------------------------------------------------
# pnpm
# ---------------------------------------------------------------------------
install_pnpm() {
  if need_cmd pnpm; then
    ok "pnpm already installed: $(pnpm --version)"
    return
  fi
  info "installing pnpm..."
  npm install -g pnpm
  ok "pnpm installed"
}

# ---------------------------------------------------------------------------
# Rust via rustup
# ---------------------------------------------------------------------------
install_rust() {
  if need_cmd rustup; then
    ok "rustup already installed: $(rustup --version 2>&1 | head -1)"
    rustup update stable --no-self-update
    rustup component add rustfmt clippy
    return
  fi
  info "installing rustup..."
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --no-modify-path
  # shellcheck source=/dev/null
  source "$HOME/.cargo/env"
  rustup component add rustfmt clippy
  ok "rust installed"
}

# ---------------------------------------------------------------------------
# Go
# ---------------------------------------------------------------------------
install_go() {
  if need_cmd go; then
    ok "go already installed: $(go version)"
    return
  fi
  if [ "$OS" = "Darwin" ]; then
    require_brew
    brew install go
  elif [ "$OS" = "Linux" ]; then
    GO_VERSION="1.22.4"
    GOARCH="amd64"
    [ "$ARCH" = "aarch64" ] && GOARCH="arm64"
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz" -o /tmp/go.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    export PATH="/usr/local/go/bin:$PATH"
  else
    die "unsupported OS: $OS"
  fi
  ok "go installed: $(go version)"
}

# ---------------------------------------------------------------------------
# Foundry (forge, cast, anvil, chisel)
# ---------------------------------------------------------------------------
install_foundry() {
  if need_cmd forge; then
    ok "foundry already installed: $(forge --version)"
    return
  fi
  info "installing foundry..."
  curl -L https://foundry.paradigm.xyz | bash
  export PATH="$HOME/.foundry/bin:$PATH"
  foundryup
  ok "foundry installed"
}

# ---------------------------------------------------------------------------
# Solana CLI
# ---------------------------------------------------------------------------
install_solana() {
  if need_cmd solana; then
    ok "solana cli already installed: $(solana --version)"
    return
  fi
  info "installing solana cli 1.18..."
  sh -c "$(curl -sSfL https://release.anza.xyz/v1.18.26/install)"
  export PATH="$HOME/.local/share/solana/install/active_release/bin:$PATH"
  ok "solana cli installed"
}

# ---------------------------------------------------------------------------
# Anchor via avm
# ---------------------------------------------------------------------------
install_anchor() {
  if need_cmd anchor; then
    ok "anchor already installed: $(anchor --version)"
    return
  fi
  if ! need_cmd avm; then
    info "installing avm (anchor version manager)..."
    cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
  fi
  info "installing anchor 0.30.1 via avm..."
  avm install 0.30.1
  avm use 0.30.1
  ok "anchor installed"
}

# ---------------------------------------------------------------------------
# Docker (check only — install platform-specific)
# ---------------------------------------------------------------------------
check_docker() {
  if need_cmd docker; then
    ok "docker already installed: $(docker --version)"
  else
    warn "docker not found — install Docker Desktop from https://docs.docker.com/desktop/"
  fi
}

# ---------------------------------------------------------------------------
# main
# ---------------------------------------------------------------------------
main() {
  info "bootstrapping titular monorepo toolchains..."
  info "os=${OS} arch=${ARCH}"

  install_mise
  install_node
  install_pnpm
  install_rust
  install_go
  install_foundry
  install_solana
  install_anchor
  check_docker

  info "installing workspace node deps..."
  pnpm install

  ok "bootstrap complete — all toolchains ready."
  ok "run: docker compose -f infra/docker/compose.dev.yml up -d"
  ok "run: pnpm -w turbo run build"
}

main "$@"
