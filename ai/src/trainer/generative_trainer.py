import random
from typing import List
from .schemas import Example

class GenerativeTrainer:
    """Huấn luyện mô hình sinh mã (ví dụ Dockerfile hoặc YAML)."""

    def __init__(self):
        self.models = {}

    def train_dockerfile_generator(self, examples: List[Example]):
        print(f"[GenerativeTrainer] Training Dockerfile generator on {len(examples)} examples...")
        patterns = [ex.output_text for ex in examples if "FROM" in ex.output_text]
        self.models["dockerfile"] = {"patterns": patterns, "trained": True}

    def train_k8s_generator(self, examples: List[Example]):
        print(f"[GenerativeTrainer] Training K8s YAML generator on {len(examples)} examples...")
        templates = [ex.output_text for ex in examples if "apiVersion" in ex.output_text]
        self.models["k8s"] = {"templates": templates, "trained": True}

    def generate(self, model_type: str) -> str:
        """Sinh ví dụ ngẫu nhiên từ mô hình."""
        if model_type not in self.models:
            return "# Model not trained yet."
        data = self.models[model_type]
        key = random.choice(list(data.keys()))
        return f"# Generated from {model_type} model using {key}"
