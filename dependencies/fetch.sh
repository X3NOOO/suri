#!/bin/bash

SELF_DIR=$(dirname $0)
cd $SELF_DIR

# Fetch piper binary (https://github.com/rhasspy/piper)
# for god's sake please do not change the naming scheme I'm not gonna do all of this again
GITHUB_API_URL="https://api.github.com/repos/rhasspy/piper/releases/latest"

detect_platform() {
  case "$(uname -s)" in
    Linux*)     PLATFORM="linux";;
    Darwin*)    PLATFORM="macos";;
    CYGWIN*|MINGW*|MSYS_NT*) PLATFORM="windows";;
    *)          PLATFORM="unknown";;
  esac
}

detect_architecture() {
  case "$(uname -m)" in
    x86_64)    ARCH="x86_64";;
    armv7l)    ARCH="armv7l";;
    aarch64)   ARCH="aarch64";;
    x86)       ARCH="x64";;
    i686)      ARCH="x64";;
    amd64)     ARCH="amd64";;
    *)         ARCH="unknown";;
  esac

  # Adjust the ARCH for windows
  if [ "$PLATFORM" = "windows" ] && [ "$ARCH" = "x64" ]; then
    ARCH="amd64"
  fi
}

detect_platform
detect_architecture


if [ "$PLATFORM" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
  echo "Unsupported platform or architecture"
  exit 1
fi

ASSET_NAME="piper_${PLATFORM}_${ARCH}.tar.gz"
ASSET_URL=$(curl -s $GITHUB_API_URL | jq -r --arg name "$ASSET_NAME" '.assets[] | select(.name == $name) | .browser_download_url')

if [ -n "$ASSET_URL" ]; then
  curl -L $ASSET_URL -o $ASSET_NAME
  echo "Download complete: $ASSET_NAME"
  
  tar -xzf $ASSET_NAME
  echo "Extraction complete."
  
  rm $ASSET_NAME
  echo "Cleanup complete: Removed $ASSET_NAME"
else
  echo "Asset not found: $ASSET_NAME"
fi
