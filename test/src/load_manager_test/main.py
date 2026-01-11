#import plotting.plot as p
from .experiment import baseline_exp, scaling_exp, selector_exp, strategy_exp, algorithm_exp, load_exp
from .plotting import plot
import numpy as np

REQUESTS = 10000

def main():
    """
    RUN EXPERIMENTS
    """
    # Base line
    bl = baseline_exp.BaselineExperiment()
    bl_result = bl.run(REQUESTS)
    print(np.mean(bl_result['result']))

    # Scaling
    scaling = scaling_exp.ScalingExperiment()
    scaling_results = scaling.run(REQUESTS)

    # Selector
    selector = selector_exp.SelectorExperiment()
    selector_results = selector.run(REQUESTS)

    # Algorithm
    algorithm = algorithm_exp.AlgorithmExperiment()
    algo_results = algorithm.run(REQUESTS)

    # Strategy
    strat = strategy_exp.StrategyExperiment()
    strat_results = strat.run(REQUESTS)


    """
    GRAPH EXPERIMENTS
    """
    plot.plot(bl_result, scaling_results, scaling.target())
    plot.plot(bl_result, selector_results, selector.target())
    plot.plot(bl_result, algo_results, algorithm.target())
    plot.plot(bl_result, strat_results, strat.target())


if __name__ == '__main__':
    main()

    
