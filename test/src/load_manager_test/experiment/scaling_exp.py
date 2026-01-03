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
            args = setup.ArgsBuilder(n=n)

            backends = args.build_backend_addr().collect_backend()
            load = ( 
                args
                .build_load_addresses()
                .build_load_queue(setup.QueueAlgorithm.FCFS)
                .build_load_selector(setup.Selector.RR)
                .build_load_strategy(setup.Strategy.MIXED)
                .collect_load()
            )

            # Execute experiment and gather result
            pids = setup.start_experiment(load_args=load, backend_args=backends)
            result = self._run_exp(num_req)
            teardown.kill_experiment(pids)

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
        
        return exper
        
