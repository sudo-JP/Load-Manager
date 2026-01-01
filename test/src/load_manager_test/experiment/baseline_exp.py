"""
Base Line Experiment is used to test on a single backend directly 
"""

from typing import override
from .exp import BaseExperience

class BaselineExperiment(BaseExperience): 
    def __init__(self): 
        super().__init__()
        self.backend_url = "http://localhost:8080/single"

    @override
    def run(self,num_req: int) -> dict:
        result = self._run_exp(num_req)

        return {
            "experiment": "Base", 
            "results": {
                'throughput': result['throughput'],
                'avg_latency': result['avg_latency'],
                'p50': result['p50'],
                'p95': result['p95'],
                'p99': result['p99'], 
                'total_time': result['total_time']
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



