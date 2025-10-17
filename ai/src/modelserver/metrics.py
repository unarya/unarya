import time
from statistics import mean

class MetricsCollector:
    """Theo dõi thời gian xử lý và số lượng yêu cầu."""

    def __init__(self):
        self.latencies = []
        self.request_count = 0

    def record_latency(self, start_time: float):
        latency = time.time() - start_time
        self.latencies.append(latency)
        self.request_count += 1

    def get_summary(self):
        if not self.latencies:
            return {"requests": 0, "avg_latency": 0.0}
        return {
            "requests": self.request_count,
            "avg_latency": round(mean(self.latencies), 4),
            "max_latency": round(max(self.latencies), 4),
            "min_latency": round(min(self.latencies), 4),
        }
