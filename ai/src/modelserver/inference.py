import hashlib
from typing import List, Dict
from .types import InferenceResult
from .cache import ResultCache
from .metrics import MetricsCollector

# Giả lập các module khác
from ai.src.preprocessor import CodePreprocessor
from ai.src.classifier.language_classifier import LanguageClassifier
from ai.src.synthesizer.dockerfile_generator import DockerfileGenerator

class InferenceEngine:
    """Bộ xử lý inference trung tâm cho server."""

    def __init__(self):
        self.cache = ResultCache(ttl=120)
        self.metrics = MetricsCollector()
        self.preprocessor = CodePreprocessor()
        self.classifier = LanguageClassifier()
        self.generator = DockerfileGenerator()

    def _cache_key(self, data: Dict) -> str:
        return hashlib.md5(str(data).encode()).hexdigest()

    def load_models(self, model_paths: List[str]):
        print(f"[InferenceEngine] Loading models: {model_paths}")
        # (Mô phỏng load, trong thực tế sẽ dùng model_loader)

    def infer(self, input_data: Dict) -> InferenceResult:
        start = time.time()
        key = self._cache_key(input_data)

        cached = self.cache.get(key)
        if cached:
            return InferenceResult(success=True, output=cached, message="Cache hit")

        op = input_data.get("operation")
        try:
            if op == "preprocess":
                code = input_data["code"]
                tokens = self.preprocessor.tokenize(code)
                encoded = self.preprocessor.encode(tokens)
                output = encoded.tolist()

            elif op == "classify":
                features = input_data["features"]
                import numpy as np
                result = self.classifier.predict(np.array(features))
                output = {"language": result.name, "confidence": result.confidence}

            elif op == "generate":
                context = input_data["context"]
                output = self.generator.generate(context)

            elif op == "health":
                output = {"status": "OK", "uptime": "active"}

            else:
                return InferenceResult(success=False, output=None, message=f"Unknown op: {op}")

            self.cache.set(key, output)
            self.metrics.record_latency(start)
            return InferenceResult(success=True, output=output)

        except Exception as e:
            return InferenceResult(success=False, output=None, message=str(e))

    def batch_infer(self, batch: List[Dict]) -> List[InferenceResult]:
        return [self.infer(item) for item in batch]
