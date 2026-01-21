import subprocess
import psycopg2
from dataclasses import dataclass
from pathlib import Path
import socket
import time

from .env import Env

from enum import Enum
from typing import Self

class QueueAlgorithm(Enum):
    FCFS = 1
    SJF = 2
    LJF = 3 
    RAND = 4
    STACK = 5

class Selector(Enum):
    RR = 1
    RANDOM = 2

class Strategy(Enum): 
    MIXED = 1
    PR = 2
    PO = 3
    PRO = 4

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

def start_backend(backend_args: list[str]) -> int:
    if not backend_args:
        raise ValueError('backend args must exist')
    backend_root = ROOT_DIR / "backend"
    if not backend_root.exists() or not backend_root.is_dir():
        raise FileNotFoundError(f"Backend directory does not exist: {backend_root}")

    proc = subprocess.Popen(backend_args, cwd=str(backend_root), start_new_session=True)
    pid = proc.pid

    # Parse host/port from args and wait until the backend is accepting TCP connections
    host = 'localhost'
    port = None
    if '--host' in backend_args:
        try:
            host = backend_args[backend_args.index('--host') + 1]
        except Exception:
            pass
    if '--port' in backend_args:
        try:
            port = int(backend_args[backend_args.index('--port') + 1])
        except Exception:
            port = None

    if port is not None:
        timeout = 10
        interval = 0.2
        start = time.time()
        while time.time() - start < timeout:
            try:
                with socket.create_connection((host, port), timeout=1):
                    break
            except Exception:
                time.sleep(interval)
        else:
            print(f'Warning: backend not responding on {host}:{port} after {timeout}s')

    return pid


def start_load_manager(load_args: list[str]) -> int:
    if not load_args:
        return -1
    load_dir = ROOT_DIR / "load-manager" / "cmd" / "load-manager"
    if not load_dir.exists() or not load_dir.is_dir():
        raise FileNotFoundError(f"Load directory does not exist: {load_dir}")

    proc = subprocess.Popen(load_args, cwd=load_dir, start_new_session=True)
    pid = proc.pid

    # Optionally wait for host:port if provided in args
    for i in range(len(load_args)): 
        if load_args[i] == '-a': 
            try:
                addr = load_args[i + 1]
                parts = addr.split(':')
                host, port = parts[0], int(parts[1])
                timeout = 10
                interval = 0.2
                start = time.time()
                while time.time() - start < timeout:
                    try:
                        with socket.create_connection((host, port), timeout=1):
                            break
                    except Exception:
                        time.sleep(interval)
                else:
                    print(f'Warning: load manager not responding on {host}:{port} after {timeout}s')
            except Exception:
                continue


    return pid

def start_experiment(load_args: list[str], backend_args: list[list[str]]) -> PIDs: 
    pids = PIDs([], -1)
    pids.backend.extend(list(map(lambda args: start_backend(args), backend_args)))
    pids.load_manager = start_load_manager(load_args)
    return pids


def reset_db() -> None:
    environment = Env()
    conn = psycopg2.connect(environment.get_db_env())
    cursor = conn.cursor()
    cursor.execute("TRUNCATE users, products, orders RESTART IDENTITY CASCADE;")
    conn.commit()
    print("Database reset!")


class Args: 
    def __init__(self, is_backend=True):

        self.args = ['go', 'run']
        if is_backend:
            self.args.append('cmd/backend/main.go')
        # Assuming this is load
        else:
            self.args.append('main.go')


    def add(self, arg: str): 
        self.args.append(arg)

GRPC_BASE_ADDR = 50050

class ArgsBuilder:
    """
    n as the number of nodes, we start our addresses at 50000
    """
    def __init__(self, n=4) -> None: 
        self.load_args = Args(is_backend=False) 
        self.backend_args = [Args() for _ in range(n)]
        self.n = n 

    """
    Build backend args
    """
    def build_backend_addr(self, host='localhost') -> Self:
        for i in range(self.n):
            self.backend_args[i].add('--host')
            self.backend_args[i].add(host)

            self.backend_args[i].add('--port')
            self.backend_args[i].add(f'{i + GRPC_BASE_ADDR}')
        return self

    def collect_backend(self) -> list[list[str]]: 
        return list(map(lambda backend: backend.args, self.backend_args))
    
    """
    Build load manager args
    """
    def build_load_queue(self, algorithm: QueueAlgorithm) -> Self: 
        self.load_args.add('-q')
        match algorithm:
            case QueueAlgorithm.FCFS: 
                self.load_args.add('FCFS')
            case QueueAlgorithm.SJF: 
                self.load_args.add('SJF')
            case QueueAlgorithm.LJF: 
                self.load_args.add('LJF')
            case QueueAlgorithm.RAND: 
                self.load_args.add('RAND')
            case QueueAlgorithm.STACK: 
                self.load_args.add('STACK')
            case _: 
                raise ValueError('Invalid Queue Algorithm')
        return self

    def build_load_selector(self, selector: Selector) -> Self: 
        self.load_args.add('-s')
        match selector: 
            case Selector.RR: 
                self.load_args.add('RR')
            case Selector.RANDOM: 
                self.load_args.add('R')
            case _: 
                raise ValueError('Invalid Selector')

        return self 
    
    def build_load_addresses(self, host='localhost') -> Self: 
        for i in range(self.n):
            self.load_args.add('-a')
            self.load_args.add(f'{host}:{i + GRPC_BASE_ADDR}')
        return self

    def build_load_strategy(self, strategy: Strategy) -> Self: 
        self.load_args.add('-l')

        match strategy: 
            case Strategy.MIXED:
                self.load_args.add('M')
            case Strategy.PO: 
                self.load_args.add('PO')
            case Strategy.PR: 
                self.load_args.add('PR')
            case Strategy.PRO: 
                self.load_args.add('PRO')
            case _: 
                raise ValueError('Invalid Strategy')
        return self

    def collect_load(self) -> list[str]: 
        return self.load_args.args
