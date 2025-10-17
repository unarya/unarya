import numpy as np
from typing import List
from .types import LanguageResult
from .model_loader import ModelLoader

class LanguageClassifier:
    """Phân loại ngôn ngữ lập trình dựa trên đặc trưng mã nguồn."""

    def __init__(self, model_loader: ModelLoader = None):
        self.loader = model_loader or ModelLoader()
        self.supported_languages = ["Python", "JavaScript", "Java", "C++", "Go", "Ruby", "PHP"]

    def predict(self, features: np.ndarray) -> LanguageResult:
        # Giả lập dự đoán bằng cách chọn ngẫu nhiên
        idx = int(features.sum()) % len(self.supported_languages)
        confidence = float((features.mean() % 1.0))
        return LanguageResult(
            name=self.supported_languages[idx],
            confidence=round(confidence, 3)
        )

    def predict_version(self, code: str) -> str:
        """Phát hiện phiên bản của ngôn ngữ (ví dụ Python 3.8, Java 11)."""
        if "async def" in code or "f\"" in code:
            return "Python 3.x"
        elif "public static void main" in code:
            return "Java 8+"
        elif "console.log" in code:
            return "JavaScript (ES6)"
        else:
            return "Unknown"
