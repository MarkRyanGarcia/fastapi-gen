import os
import sys

__version__ = os.environ.get("FAPI_INIT_VERSION", "0.0.0")


def main() -> None:
    from ._binary import ensure_binary
    import subprocess

    binary = ensure_binary(__version__)
    result = subprocess.run([str(binary)] + sys.argv[1:])
    sys.exit(result.returncode)
