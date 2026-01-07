from .setup import PIDs 
import signal
import os 

def kill_process(pid: int): 
    print(f'Killing {pid}')
    os.kill(pid, signal.SIGINT)

"""
Not my fault the naming convention is like this 
"""
def kill_experiment(pids: PIDs): 
    map(lambda pid: kill_process(pid), pids.backend)
    os.kill(pids.load_manager, signal.SIGINT)
