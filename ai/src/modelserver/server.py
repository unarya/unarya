import grpc
from concurrent import futures
import time

from ..pb import ai_pb2, ai_pb2_grpc
from .inference import InferenceEngine


class AIServiceServicer(ai_pb2_grpc.AIServiceServicer):
    def __init__(self):
        self.engine = InferenceEngine()

    def AnalyzeCode(self, request, context):
        print(f"[AIService] Received AnalyzeCode for language={request.language}")

        result = self.engine.infer({
            "operation": "AnalyzeCode",
            "language": request.language,
            "code_structure": request.code_structure
        })

        # Check success first
        if not result.get("success"):
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(result.get("message", "Unknown error"))
            return ai_pb2.AIAnalyzeResponse()

        # Get actual output data
        output = result.get("output", {})

        response = ai_pb2.AIAnalyzeResponse(
            insights=output.get("insights", "Unknown"),
            confidence=str(output.get("confidence", "0.0"))
        )
        return response


class ModelServer:
    def __init__(self, port: int = 6000):
        self.port = port
        self.server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
        ai_pb2_grpc.add_AIServiceServicer_to_server(AIServiceServicer(), self.server)
        self.server.add_insecure_port(f"[::]:{self.port}")

    def serve(self):
        print(f"[ModelServer] 🚀 Starting gRPC server on port {self.port} ...")
        self.server.start()
        print("[ModelServer] Available RPCs: AnalyzeCode()")

        try:
            while True:
                time.sleep(60)
        except KeyboardInterrupt:
            print("\n[ModelServer] 🛑 Stopping server...")
            self.server.stop(0)


if __name__ == "__main__":
    server = ModelServer(port=6000)
    server.serve()
