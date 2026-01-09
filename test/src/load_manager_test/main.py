#import plotting.plot as p
from load_manager_test.experiment import baseline_exp, scaling_exp, selector_exp, strategy_exp, algorithm_exp, load_exp
from load_manager_test.plotting import plot

REQUESTS = 2

def main():
    """
    RUN EXPERIMENTS
    """
    # Base line
    bl = baseline_exp.BaselineExperiment()
    bl_result = bl.run(REQUESTS)

    # Scaling
    scaling = scaling_exp.ScalingExperiment()
    scaling_results = bl.run(REQUESTS)

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

    
