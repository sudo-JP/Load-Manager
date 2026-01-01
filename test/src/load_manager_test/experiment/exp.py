from abc import ABC, abstractmethod
import time 
import requests 

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

    def _run_exp(self, num_req: int) -> dict: 
        print(f"Running {num_req} requests to backend...")

        # Number of requests
        for i in range(num_req): 
            user_data = {
                "name": f'Test User {i}', 
                "email": f'user{i}@example.com', 
                'password': f'user{i}somethingsomething',
            }
            
            start_time = time.perf_counter()
            try:
                requests.post(f'{self.backend_url}/{self.userRoute}',
                                json=user_data,
                                timeout=5)

                end_time = time.perf_counter()
                elapsed = (end_time - start_time) * 1000
                self.latencies.append(elapsed)
            except Exception as e: 
                print(f'Request {i} failed {e}')

        if not self.latencies:
            raise ValueError("No succcessful req")

        # Calculation

        sorted_lat = sorted(self.latencies) 
        total_time = sum(self.latencies) 
        avg = sum(sorted_lat) / len(sorted_lat)
        p50 = sorted_lat[len(sorted_lat) // 2]
        p95 = sorted_lat[int(len(sorted_lat) * 0.95)]
        p99 = sorted_lat[int(len(sorted_lat) * 0.99)]
        return {
            'throughput': num_req / total_time,
            'avg_latency': avg, 
            'p50': p50, 
            'p95': p95, 
            'p99': p99, 
            'total_time': total_time
        }
