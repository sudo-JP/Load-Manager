#import plotting.plot as p
from load_manager_test.experiment import baseline_exp, scaling_exp
from load_manager_test.plotting import plot

REQUESTS = 2

if __name__ == '__main__':
    bl = baseline_exp.BaselineExperiment()
    bl_result = bl.run(REQUESTS)

    scaling = scaling_exp.ScalingExperiment()
    scaling_results = bl.run(REQUESTS)

    plot.plot(bl_result, scaling_results, scaling.target())

#start_backend(1)

    
