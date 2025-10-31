import numpy as np
from .schemas import Dataset, EvalResult

class Evaluator:
    """Đánh giá mô hình sau huấn luyện."""

    def evaluate(self, predictions: np.ndarray, targets: np.ndarray) -> EvalResult:
        if predictions.shape != targets.shape:
            raise ValueError("Predictions and targets must have the same shape")

        accuracy = float((predictions == targets).mean())
        loss = float(((predictions - targets) ** 2).mean())
        return EvalResult(accuracy=accuracy, loss=loss)

    def evaluate_dataset(self, dataset: Dataset) -> EvalResult:
        """Mô phỏng đánh giá dựa trên số mẫu và độ dài trung bình."""
        avg_len = np.mean([len(ex.input_text) for ex in dataset.examples])
        score = 1 / (1 + np.exp(-avg_len / 100))
        return EvalResult(accuracy=round(score, 3), loss=round(1 - score, 3))
