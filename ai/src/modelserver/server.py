import grpc
from concurrent import futures
import time
from typing import Any
from .types import InferenceRequest, Response
from .inference import InferenceEngine

# Giả lập thay cho proto real
class ModelServer:
    """gRPC server phục vụ inference."""

    def __init__(self):
        self.engine = InferenceEngine()
        self.server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
        self.running = False

    def handle_request(self, request: InferenceRequest) -> Response:
        print(f"[ModelServer] Handling {request.operation}")
        result = self.engine.infer({"operation": request.operation, **request.payload})
        status = "OK" if result.success else "ERROR"
        return Response(status=status, result=result)

    def serve(self, port: int = 50051):
        """Chạy server gRPC (mô phỏng, không cần .proto)."""
        self.running = True
        print(f"[ModelServer] Serving on port {port} ...")
        print("Available endpoints: Preprocess, Classify, Generate, HealthCheck")
        try:
            while self.running:
                time.sleep(2)
        except KeyboardInterrupt:
            print("\n[ModelServer] Stopping server...")
            self.running = False

    def stop(self):
        self.running = False
