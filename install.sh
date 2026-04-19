#!/bin/sh
set -e

REPO="moriT958/libpukiwiki"
BINARY="md2pw"
INSTALL_DIR="${HOME}/.local/bin"

# OS と Arch を検出
OS=$(uname -s)
ARCH=$(uname -m)

# goreleaser の命名規則に合わせる
case "$ARCH" in
  x86_64) ARCH="x86_64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  i386|i686) ARCH="i386" ;;
esac

# 最新バージョンを取得
VERSION=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

# ダウンロードURLを構築
URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}_${OS}_${ARCH}.tar.gz"

# インストールディレクトリを作成
mkdir -p "$INSTALL_DIR"

echo "Downloading ${BINARY} ${VERSION}..."
curl -sSfL "$URL" | tar xz -C /tmp

echo "Installing to ${INSTALL_DIR}..."
mv "/tmp/${BINARY}" "${INSTALL_DIR}/"
chmod +x "${INSTALL_DIR}/${BINARY}"

# PATH 設定
add_to_path() {
  SHELL_NAME=$(basename "$SHELL")
  case "$SHELL_NAME" in
    bash)
      RC_FILE="$HOME/.bashrc"
      ;;
    zsh)
      RC_FILE="$HOME/.zshrc"
      ;;
    *)
      RC_FILE=""
      ;;
  esac

  if [ -n "$RC_FILE" ] && [ -f "$RC_FILE" ]; then
    if ! grep -q "${INSTALL_DIR}" "$RC_FILE" 2>/dev/null; then
      echo "" >> "$RC_FILE"
      echo "# md2pw" >> "$RC_FILE"
      echo "export PATH=\"\$PATH:${INSTALL_DIR}\"" >> "$RC_FILE"
      echo "Added ${INSTALL_DIR} to PATH in ${RC_FILE}"
      echo "Run 'source ${RC_FILE}' or restart your terminal"
    fi
  fi
}

# PATH にインストールディレクトリがなければ追加
case ":$PATH:" in
  *":${INSTALL_DIR}:"*)
    # already in PATH
    ;;
  *)
    add_to_path
    ;;
esac

echo "Done! Run 'md2pw -o <outfile> <input.md>' to get started."
