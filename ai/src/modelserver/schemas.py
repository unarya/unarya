from dataclasses import dataclass
from typing import Dict, List, Any, Optional

@dataclass
class InferenceRequest:
    """Yêu cầu inference từ client."""
    operation: str  # preprocess | classify | generate | health
    payload: Dict[str, Any]

@dataclass
class InferenceResult:
    """Kết quả inference."""
    success: bool
    output: Any
    message: Optional[str] = None

@dataclass
class Response:
    """Phản hồi trả về cho client."""
    status: str
    result: Optional[InferenceResult] = None

@dataclass
class CacheEntry:
    """Mục trong cache."""
    key: str
    value: Any
    timestamp: float
