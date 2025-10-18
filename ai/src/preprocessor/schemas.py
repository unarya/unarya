from dataclasses import dataclass
from typing import Dict

@dataclass
class Token:
    """Đại diện cho một token mã nguồn."""
    type: str
    value: str
    position: int
    metadata: Dict
