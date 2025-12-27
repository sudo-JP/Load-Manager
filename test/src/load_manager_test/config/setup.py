import subprocess
from dataclasses import dataclass
from pathlib import Path

"""
Set up root dir, travel .. 5 times
"""
ROOT_DIR = Path(__file__).resolve()

for _ in range(5): 
    ROOT_DIR = ROOT_DIR.parent

@dataclass
class PIDs: 
    backend: list[int]
    load_manager: int

def start_backend(n: int, backend_args: list[str]) -> list[int]: 
    backend_dir = ROOT_DIR / "backend" / "cmd" / "backend" 
    if not backend_dir.exists() or not backend_dir.is_dir():
        raise FileNotFoundError(f"Backend directory does not exist: {backend_dir}")

    pids = []
    for _ in range(n): 
        pids.append(subprocess.Popen(backend_args, cwd=backend_dir).pid)
    return pids

def start_load_manager(load_args: list[str]) -> int:
    load_dir = ROOT_DIR / "load-manager" / "cmd" / "load-manager"
    if not load_dir.exists() or not load_dir.is_dir():
        raise FileNotFoundError(f"Load directory does not exist: {load_dir}")

    return subprocess.Popen(load_args, cwd=load_dir).pid

"""
n is the number of backend nodes, if n == 1, then load manager won't be started
"""
def start_process(n: int, load_args: list[str], backend_args: list[str]) -> PIDs: 
    pids = PIDs([], -1)
    pids.backend = start_backend(n, backend_args)
    if n != 1: 
        pids.load_manager = start_load_manager(load_args)
    return pids

