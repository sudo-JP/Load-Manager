"""
Purpose: Does adding backends improve performance? 
"""


from typing import override
import time 
import requests 
from .exp import BaseExperience

class NumberNodesExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()

    @override
    def run(self, num_req: int) -> dict:

        exper = {
            "experiment": "Scaling",
            "results": []
        }

        # Number of backend nodes
        for n in [2, 4, 8, 16]:
            
            # Number of requests
            for i in range(num_req): 
                user_data = {
                    "name": f'Test User {i}', 
                    "email": f'user{i}@example.com', 
                    'password': f'user{i}somethingsomething',
                }
                
                start_time = time.perf_counter()
                try:
                    requests.post(f'{self.backend_url}/balancer/{self.userRoute}',
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


            exper["results"].append({
                'throughput': num_req / total_time,
                'avg_latency': avg, 
                'p50': p50, 
                'p95': p95, 
                'p99': p99, 
                'total_time': total_time
            })
        
        return exper
        
