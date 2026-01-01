"""
Purpose: Which strategy is the best? 
"""

from typing import override
from .exp import BaseExperience

class StrategyExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()

    @override
    def run(self, num_req: int) -> dict:
        exper = {
            "experiment": "Strategy",
            "results": []
        }

        # Strategy
        for strat in ['M', 'PR', 'PO', 'PRO']:

            # TODO: Call config set up 
            
            result = self._run_exp(num_req)

            exper["results"].append({
                'nodes': 4, 
                'algorithm': 'FCFS',
                'selector': 'RR', 
                'strategy': strat,
                'throughput': result['throughput'],
                'avg_latency': result['avg_latency'],
                'p50': result['p50'],
                'p95': result['p95'],
                'p99': result['p99'], 
                'total_time': result['total_time']
            })
        
        return exper
        
