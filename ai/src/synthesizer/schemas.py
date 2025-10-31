from dataclasses import dataclass
from typing import List, Dict, Optional

@dataclass
class CodeContext:
    """Thông tin ngữ cảnh của mã nguồn để sinh Dockerfile."""
    language: str
    dependencies: List[str]
    entry_point: str
    version: Optional[str] = None

@dataclass
class Service:
    """Mô tả một service trong docker-compose.yml."""
    name: str
    image: str
    ports: List[str]
    environment: Dict[str, str] = None
    volumes: List[str] = None
    depends_on: List[str] = None

@dataclass
class AppConfig:
    """Cấu hình ứng dụng cho K8s deployment."""
    name: str
    image: str
    replicas: int = 1
    ports: List[int] = None
    env: Dict[str, str] = None
    domain: Optional[str] = None

@dataclass
class ValidationResult:
    """Kết quả kiểm tra file sinh ra."""
    valid: bool
    errors: List[str]
