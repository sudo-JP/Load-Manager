"""
Purpose: Which strategy is the best? 
"""

from typing import override
from load_manager_test.experiment.exp import BaseExperience
from load_manager_test.config import setup, teardown

class StrategyExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()

    @override
    def target(self) -> str: 
        return 'strategy'

    @override
    def run(self, num_req: int) -> dict:
        exper = {
            "experiment": "Strategy",
            "results": []
        }

        # Strategy
        nodes = 4
        for strat in [setup.Strategy.MIXED, setup.Strategy.PO, setup.Strategy.PR, setup.Strategy.PRO]:
            setup.reset_db()
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
                'result': result
            })
        
        return exper
        
