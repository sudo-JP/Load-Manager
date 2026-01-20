"""
Purpose: Which queue algorithm is the best? 
"""

from typing import override
from .exp import BaseExperience
from ..config import setup, teardown

class AlgorithmExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()

    @override
    def target(self) -> str: 
        return 'algorithm'

    @override
    def run(self, num_req: int) -> dict:
        exper = {
            "experiment": "Algorithm",
            "results": []
        }

        # Algorithm
        nodes = 4
        for algo in [setup.QueueAlgorithm.FCFS, setup.QueueAlgorithm.SJF, setup.QueueAlgorithm.LJF, setup.QueueAlgorithm.RAND]:
            
            args = setup.ArgsBuilder(n=nodes)

            backends = args.build_backend_addr().collect_backend()
            load = ( 
                args
                .build_load_addresses()
                .build_load_queue(algo)
                .build_load_selector(setup.Selector.RR)
                .build_load_strategy(setup.Strategy.MIXED)
                .collect_load()
            )

            # Execute experiment and gather result
            pids = setup.start_experiment(load_args=load, backend_args=backends)
            result = self._run_exp(num_req)
            teardown.kill_experiment(pids)

            exper["results"].append({
                'nodes': nodes, 
                'algorithm': algo.name,
                'selector': 'RR',
                'strategy': 'M',
                'result': result
            })
        
        return exper
        
