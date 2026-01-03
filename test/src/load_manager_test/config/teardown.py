from .setup import PIDs 
import signal
import os 

"""
Not my fault the naming convention is like this 
"""
def kill_experiment(pids: PIDs): 
    map(lambda pid: os.kill(pid, signal.SIGINT), pids.backend)
    os.kill(pids.load_manager, signal.SIGINT)
