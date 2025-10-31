from dataclasses import dataclass
from typing import Any, List

@dataclass
class Config:
    """Cấu hình huấn luyện mô hình."""
    learning_rate: float = 1e-4
    batch_size: int = 32
    epochs: int = 5
    embedding_dim: int = 128
    model_name: str = "default_model"

@dataclass
class Example:
    """Một ví dụ huấn luyện (ví dụ mã và đầu ra mong đợi)."""
    input_text: str
    output_text: str

@dataclass
class Dataset:
    """Tập dữ liệu huấn luyện."""
    examples: List[Example]

@dataclass
class EvalResult:
    """Kết quả đánh giá mô hình."""
    accuracy: float
    loss: float
    details: Any = None
