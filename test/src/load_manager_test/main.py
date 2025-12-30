#import plotting.plot as p
from config.setup import start_backend
from experiment.baseline_exp import BaselineExperiment

NUM_REQS = 100
BACKEND_NODES = 4


if __name__ == '__main__':
    bl = BaselineExperiment(100)
    bl.run()
#start_backend(1)

    
