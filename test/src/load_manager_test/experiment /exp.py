from abc import ABC, abstractmethod

class BaseExperience(ABC): 

    def __init__(self, n: int):
        self.results = {}
        self.num_req = n
        self.backend_url = "https://localhost:9000"
        self.latencies = []
        self.userRoute = "user"
        self.productRoute = "product"
        self.orderRoute = "order"

    @abstractmethod
    def run(self): 
        "Execute experiment"
        pass

    def collect_results(self) -> dict:
        return self.results
