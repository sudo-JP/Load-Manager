#import plotting.plot as p
from .config.setup import start_backend
from .experiment.scaling_exp import ScalingExperiment

NUM_REQS = 100
BACKEND_NODES = 4


if __name__ == '__main__':
    bl = ScalingExperiment()
    bl.run(1000)
#start_backend(1)

    
