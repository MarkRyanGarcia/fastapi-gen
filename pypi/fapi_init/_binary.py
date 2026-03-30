"""Downloads and caches the fapi-init binary from GitHub Releases."""

import os
import platform
import stat
import sys
import tarfile
import urllib.request
import zipfile
from pathlib import Path

GITHUB_REPO = "markryangarcia/fapi-init"
BIN_NAME = "fapi-init"
CACHE_DIR = Path.home() / ".cache" / "fapi-init"


def _platform_asset() -> str:
    system = platform.system().lower()
    machine = platform.machine().lower()

    arch = "arm64" if machine in ("arm64", "aarch64") else "amd64"

    if system == "darwin":
        return f"fapi-init_darwin_{arch}.tar.gz"
    elif system == "linux":
        return f"fapi-init_linux_{arch}.tar.gz"
    elif system == "windows":
        return f"fapi-init_windows_{arch}.zip"
    else:
        raise RuntimeError(f"Unsupported platform: {system}/{machine}")


def _bin_path(version: str) -> Path:
    suffix = ".exe" if platform.system().lower() == "windows" else ""
    return CACHE_DIR / version / (BIN_NAME + suffix)


def ensure_binary(version: str) -> Path:
    bin_path = _bin_path(version)
    if bin_path.exists():
        return bin_path

    asset = _platform_asset()
    url = f"https://github.com/{GITHUB_REPO}/releases/download/v{version}/{asset}"

    bin_path.parent.mkdir(parents=True, exist_ok=True)

    archive_path = bin_path.parent / asset
    print(f"Downloading fapi-init v{version}...", file=sys.stderr)
    urllib.request.urlretrieve(url, archive_path)

    if asset.endswith(".zip"):
        with zipfile.ZipFile(archive_path) as zf:
            zf.extractall(bin_path.parent)
    else:
        with tarfile.open(archive_path, "r:gz") as tf:
            tf.extractall(bin_path.parent)

    archive_path.unlink()

    # ensure executable
    bin_path.chmod(bin_path.stat().st_mode | stat.S_IEXEC | stat.S_IXGRP | stat.S_IXOTH)
    return bin_path
