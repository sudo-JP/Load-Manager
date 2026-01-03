"""
Purpose: How does system handle increasing load? 
"""

from typing import override
from .exp import BaseExperience

class LoadExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()

    @override
    def run(self, num_req: int) -> dict:
        # TODO: Create the base line classs and compare that against some fixed variable 
        return {}
        
