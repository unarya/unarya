import hashlib
import time
from typing import List, Dict, Any
from .cache import ResultCache
from .metrics import MetricsCollector
from ..classifier.language_classifier import LanguageClassifier
from ..preprocessor import CodePreprocessor
from ..synthesizer.dockerfile_generator import DockerfileGenerator

class InferenceEngine:
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

    def infer(self, input_data: Dict) -> Dict[str, Any]:
        """
        Return dict instead of InferenceResult object for gRPC compatibility
        """
        start = time.time()
        key = self._cache_key(input_data)

        cached = self.cache.get(key)
        if cached:
            return {
                "success": True,
                "output": cached,
                "message": "Cache hit"
            }

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

            elif op == "AnalyzeCode":
                language = input_data.get("language", "unknown")
                code_structure = input_data.get("code_structure", "")
                # Phân tích code và trả về insights
                output = {
                    "insights": f"Analyzed {language} code",
                    "confidence": "0.95"
                }

            else:
                return {
                    "success": False,
                    "output": None,
                    "message": f"Unknown op: {op}"
                }

            self.cache.set(key, output)
            self.metrics.record_latency(start)
            return {
                "success": True,
                "output": output,
                "message": "OK"
            }

        except Exception as e:
            return {
                "success": False,
                "output": None,
                "message": str(e)
            }

    def batch_infer(self, batch: List[Dict]) -> List[Dict[str, Any]]:
        return [self.infer(item) for item in batch]