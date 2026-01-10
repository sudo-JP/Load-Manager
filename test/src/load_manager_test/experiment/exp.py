from abc import ABC, abstractmethod
import time 
import requests 
import numpy as np
from ..config import setup

class BaseExperience(ABC): 

    def __init__(self):
        self.results = {}
        self.backend_url = "http://localhost:8000/balancer"
        self.userRoute = "user"
        self.productRoute = "product"
        self.orderRoute = "order"

    @abstractmethod
    def run(self, num_req: int) -> dict:
        "Execute experiment"
        pass

    @abstractmethod
    def target(self) -> str: 
        "Target key for the experiment"
        pass

    def _run_exp(self, num_req: int) -> np.ndarray:
        setup.reset_db()
        print(f"Running {num_req} requests to backend...")
        print(f'{self.backend_url}/{self.userRoute}')

        latencies = np.empty(num_req)
        success_count = 0

        # Number of requests
        for i in range(num_req):
            user_data = {
                "name": f'Test User {i}',
                "email": f'user{i}@example.com',
                'password': f'user{i}somethingsomething',
            }

            start_time = time.perf_counter()
            try:
                resp = requests.post(f'{self.backend_url}/{self.userRoute}',
                                json=user_data,
                                timeout=5)

                if resp.status_code >= 400:
                    print(f'Request {i} failed status {resp.status_code}: {resp.text}')
                else:
                    end_time = time.perf_counter()
                    elapsed = (end_time - start_time) * 1000
                    latencies[i] = (elapsed)
                    success_count += 1
            except Exception as e:
                print(f'Request {i} failed {e}')

        if success_count == 0:
            raise ValueError("No succcessful req")

        # Calculation

        """sorted_lat = sorted(self.latencies) 
        total_time = sum(self.latencies) 
        avg = sum(sorted_lat) / len(sorted_lat)
        p50 = sorted_lat[len(sorted_lat) // 2]
        p95 = sorted_lat[int(len(sorted_lat) * 0.95)]
        p99 = sorted_lat[int(len(sorted_lat) * 0.99)]"""
        return latencies
