from typing import List, Dict
import numpy as np
from .tokenizer import CodeTokenizer
from .encoder import CodeEncoder
from .tree_encoder import TreeEncoder
from .normalizer import Normalizer
from .types import Token

class CodePreprocessor:
    """Pipeline tiền xử lý mã nguồn."""

    def __init__(self):
        self.tokenizer = CodeTokenizer()
        self.encoder = CodeEncoder()
        self.tree_encoder = TreeEncoder()
        self.normalizer = Normalizer()

    def tokenize(self, code: str) -> List[Token]:
        return self.tokenizer.tokenize(code)

    def encode(self, tokens: List[Token]) -> np.ndarray:
        return self.encoder.encode(tokens)

    def encode_tree(self, ast: Dict) -> np.ndarray:
        return self.tree_encoder.encode_tree(ast)

    def normalize(self, matrix: np.ndarray) -> np.ndarray:
        return self.normalizer.normalize(matrix)
