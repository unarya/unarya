from dataclasses import dataclass
from typing import List, Optional

@dataclass
class LanguageResult:
    """Kết quả phân loại ngôn ngữ lập trình."""
    name: str
    confidence: float

@dataclass
class Framework:
    """Thông tin framework được phát hiện."""
    name: str
    confidence: float

@dataclass
class StyleMetrics:
    """Phân tích phong cách lập trình."""
    indentation_spaces: int
    naming_convention: str  # snake_case, camelCase, PascalCase
    average_line_length: float
    comment_ratio: float
    complexity_estimate: float
