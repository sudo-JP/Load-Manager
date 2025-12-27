from typing import override
import time 
import requests as req
from .exp import BaseExperience

class BaselineExperiment(BaseExperience): 
    def __init__(self, n: int): 
        super().__init__(n)

    @override
    def run(self):
        start_time = time.perf_counter()
        for _ in range(self.n): 
            req.get(self.userRoute)

        end_time = time.perf_counter()
        elapsed = end_time - start_time

        self.results['time'] = elapsed


