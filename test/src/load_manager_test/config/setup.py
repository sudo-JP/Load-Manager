import subprocess
from pathlib import Path

"""
Set up root dir, travel .. 5 times
"""
ROOT_DIR = Path(__file__).resolve()

for _ in range(5): 
    ROOT_DIR = ROOT_DIR.parent


def start_backend(n: int, backend_args: list[str]) -> None: 
    backend_dir = ROOT_DIR / "backend" / "cmd" / "backend" 
    if not backend_dir.exists() or not backend_dir.is_dir():
        raise FileNotFoundError(f"Backend directory does not exist: {backend_dir}")

    for _ in range(n): 
        subprocess.run(backend_args, cwd=backend_dir)

def start_load_manager(load_args: list[str]):
    load_dir = ROOT_DIR / "load-manager" / "cmd" / "load-manager"
    if not load_dir.exists() or not load_dir.is_dir():
        raise FileNotFoundError(f"Load directory does not exist: {load_dir}")

    subprocess.run(load_args, cwd=load_dir)

"""
n is the number of backend nodes, if n == 1, then load manager won't be started
"""
def start_process(n: int, load_args: list[str], backend_args: list[str]) -> None: 
    start_backend(n, backend_args)
    if n != 1: 
        start_load_manager(load_args)

