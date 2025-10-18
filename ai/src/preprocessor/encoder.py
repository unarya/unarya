import numpy as np
from typing import List
from .schemas import Token

class CodeEncoder:
    """Chuyển danh sách token thành vector số."""

    def __init__(self):
        self.vocab = {}
        self.next_index = 1

    def _get_index(self, token: str) -> int:
        if token not in self.vocab:
            self.vocab[token] = self.next_index
            self.next_index += 1
        return self.vocab[token]

    def encode(self, tokens: List[Token]) -> np.ndarray:
        """Encode token thành vector chỉ số số học."""
        indices = [self._get_index(tok.value) for tok in tokens]
        return np.array(indices, dtype=np.float32)
