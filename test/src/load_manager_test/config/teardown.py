from .setup import PIDs
import signal
import os 

def kill_process(pid: int):
    print(f'Killing {pid}')
    try:
        # Try killing the whole process group if the process was started in its own session
        os.killpg(pid, signal.SIGINT)
    except Exception:
        try:
            os.kill(pid, signal.SIGINT)
        except Exception as e:
            print(f'Failed to kill {pid}: {e}')

"""
Not my fault the naming convention is like this 
"""

def kill_experiment(pids: PIDs):
    for pid in pids.backend:
        kill_process(pid)
    try:
        os.killpg(pids.load_manager, signal.SIGINT)
    except Exception:
        try:
            os.kill(pids.load_manager, signal.SIGINT)
        except Exception as e:
            print(f'Failed to kill load_manager {pids.load_manager}: {e}')
