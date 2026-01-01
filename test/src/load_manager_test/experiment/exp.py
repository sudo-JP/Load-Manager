from abc import ABC, abstractmethod

class BaseExperience(ABC): 

    def __init__(self):
        self.results = {}
        self.backend_url = "http://localhost:9000"
        self.latencies = []
        self.userRoute = "user"
        self.productRoute = "product"
        self.orderRoute = "order"

    @abstractmethod
    def run(self, num_req: int) -> dict:
        "Execute experiment"
        pass

    def collect_results(self) -> dict:
        return self.results
