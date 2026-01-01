"""
Base Line Experiment is used to test on a single backend directly 
"""

from typing import override
import time 
import requests 
from .exp import BaseExperience

class BaselineExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()
        self.backend_url = "http://localhost:808"

    @override
    def run(self,num_req: int) -> dict:
        print(f"Running {num_req} requests to backend...")


        # Just POST for now 
        for i in range(num_req): 
            user_data = {
                "name": f"Test User {i}",
                "email": f"user{i}@example.com", 
                "password": f"user{i}somethingsomething",
            }

            start_time = time.perf_counter()
            try: 
                requests.post(f"{self.backend_url}/single/user",
                                     json=user_data, 
                                     timeout=5)
                end_time = time.perf_counter()
                elapsed = (end_time - start_time) * 1000
                self.latencies.append(elapsed)
            except Exception as e: 
                print(f"Request {i} failed {e}")

        if not self.latencies:
            raise ValueError("No succcessful req")

        sorted_lat = sorted(self.latencies)
        avg = sum(sorted_lat) / len(sorted_lat)
        p50 = sorted_lat[len(sorted_lat) // 2]
        p95 = sorted_lat[int(len(sorted_lat) * 0.95)]
        p99 = sorted_lat[int(len(sorted_lat) * 0.99)]

        total_time = sum(self.latencies)

        return {
            "experiment": "base line", 
            "results": {
                'throughput': num_req / total_time,
                'avg_latency': avg, 
                'p50': p50, 
                'p95': p95, 
                'p99': p99, 
                'total_time': total_time
            }
        }
        
    def print_results(self): 
        if not self.latencies:
            print("No succcessful req")
            return 

        sorted_lat = sorted(self.latencies)
        avg = sum(sorted_lat) / len(sorted_lat)
        p50 = sorted_lat[len(sorted_lat) // 2]
        p95 = sorted_lat[int(len(sorted_lat) * 0.95)]
        p99 = sorted_lat[int(len(sorted_lat) * 0.99)]

        print("RESULTS")
        print(f"Total req: {len(self.latencies)}")
        print(f"Avg latency: {avg:.2f}ms")
        print(f"P50: {p50:.2f}ms")
        print(f"P95: {p95:.2f}ms")
        print(f"P99: {p99:.2f}ms")
        print(f"Min: {min(sorted_lat):.2f}ms")
        print(f"Max: {max(sorted_lat):.2f}ms")

if __name__ == '__main__':
    exp = BaselineExperiment()
    exp.run(1000)



