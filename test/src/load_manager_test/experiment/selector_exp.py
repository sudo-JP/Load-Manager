"""
Purpose: Which selector is the best? 
"""

from typing import override
from .exp import BaseExperience
from ..config import setup

class SelectorExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()

    @override
    def run(self, num_req: int) -> dict:
        exper = {
            "experiment": "Selector",
            "results": []
        }

        # Selector
        nodes = 4
        for selector in ['RR', 'Random']:

            # TODO: Call config set up 
            addrs = setup.generate_address(nodes)

            # Fix
            addr_args = list(map(lambda addr: f'-a {addr}', addrs))
            args = setup.default_arg()
            args.extend(addr_args)
            args.extend(['-q', 'fcfs', '-s', selector, '-l', 'M'])

            
            result = self._run_exp(num_req)

            exper["results"].append({
                'nodes': nodes, 
                'algorithm': 'FCFS',
                'selector': selector, 
                'strategy': 'M',
                'throughput': result['throughput'],
                'avg_latency': result['avg_latency'],
                'p50': result['p50'],
                'p95': result['p95'],
                'p99': result['p99'], 
                'total_time': result['total_time']
            })
        
        return exper
        
