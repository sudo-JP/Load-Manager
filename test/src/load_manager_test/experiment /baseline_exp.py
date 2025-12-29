from typing import override
import time 
import requests 
#from .exp import BaseExperience

class BaselineExperiment(): 
    def __init__(self, n: int): 
        self.results = {}
        self.num_req = n
        self.backend_url = "http://localhost:9000"
        self.latencies = []
        self.userRoute = "user"
        self.productRoute = "product"
        self.orderRoute = "order"

    def run(self):
        print(f"Running {self.num_req} requests to backend...")

        # Just POST for now 
        for i in range(self.num_req): 
            user_data = {
                "name": f"Test User {i}",
                "email": f"user{i}@example.com", 
                "password": f"user{i}",
            }

            start_time = time.perf_counter()
            try: 
                requests.post(f"{self.backend_url}/single/{self.userRoute}",
                                     json=user_data, 
                                     timeout=5)
                end_time = time.perf_counter()
                elapsed = (end_time - start_time) * 1000
                self.latencies.append(elapsed)

                if i % 100 == 0: 
                    print(f"Completed {i}/{self.num_req}")
            except Exception as e: 
                print(f"Request {i} failed {e}")

        self.print_results()

        
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
    exp = BaselineExperiment(100)
    exp.run()



