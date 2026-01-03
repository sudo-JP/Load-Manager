"""
Purpose: Which strategy is the best? 
"""

from typing import override
from .exp import BaseExperience
from ..config import setup, teardown

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
        nodes = 4
        for strat in [setup.Strategy.MIXED, setup.Strategy.PO, setup.Strategy.PR, setup.Strategy.PRO]:
            args = setup.ArgsBuilder(n=nodes)

            backends = args.build_backend_addr().collect_backend()
            load = ( 
                args
                .build_load_addresses()
                .build_load_queue(setup.QueueAlgorithm.FCFS)
                .build_load_selector(setup.Selector.RR)
                .build_load_strategy(strat)
                .collect_load()
            )

            # Execute experiment and gather result
            pids = setup.start_experiment(load_args=load, backend_args=backends)
            result = self._run_exp(num_req)
            teardown.kill_experiment(pids)

            exper["results"].append({
                'nodes': nodes, 
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
        
