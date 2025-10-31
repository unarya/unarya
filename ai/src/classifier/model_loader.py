import joblib
import os
import numpy as np

class ModelLoader:
    """Tải và cache các mô hình phân loại được huấn luyện sẵn."""

    def __init__(self, model_dir: str = "models/"):
        self.model_dir = model_dir
        self.cache = {}

    def load_model(self, name: str):
        """Tải mô hình từ file, nếu có cache thì dùng lại."""
        if name in self.cache:
            return self.cache[name]

        path = os.path.join(self.model_dir, f"{name}.pkl")
        if not os.path.exists(path):
            raise FileNotFoundError(f"Model file not found: {path}")

        model = joblib.load(path)
        self.cache[name] = model
        return model

    def predict(self, model_name: str, features: np.ndarray):
        model = self.load_model(model_name)
        return model.predict(features)
