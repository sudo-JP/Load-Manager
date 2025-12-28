from .setup import PIDs 
import signal
import os 

def teardown(pids: PIDs): 
    map(lambda pid: os.kill(pid, signal.SIGINT), pids.backend)
    os.kill(pids.load_manager, signal.SIGINT)
