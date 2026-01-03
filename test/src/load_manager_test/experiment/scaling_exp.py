"""
Purpose: Does adding backends improve performance? 
"""

from typing import override
from .exp import BaseExperience
from ..config import setup, teardown

class NumberNodesExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()

    @override
    def run(self, num_req: int) -> dict:
        exper = {
            "experiment": "Scaling",
            "results": []
        }

        # Number of backend nodes
        for n in [2, 4, 8, 16]:

            # TODO: Call config set up 
            pids = setup.start_process(n, [], [])
            
            result = self._run_exp(num_req)

            exper["results"].append({
                'nodes': n, 
                'algorithm': 'FCFS',
                'selector': 'RR', 
                'strategy': 'M',
                'throughput': result['throughput'],
                'avg_latency': result['avg_latency'],
                'p50': result['p50'],
                'p95': result['p95'],
                'p99': result['p99'], 
                'total_time': result['total_time']
            })

            teardown.teardown(pids)
        
        return exper
        
