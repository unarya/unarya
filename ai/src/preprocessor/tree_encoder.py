import numpy as np
from typing import Dict, Any

class TreeEncoder:
    """Mã hóa cây cú pháp trừu tượng (AST) thành vector."""

    def __init__(self):
        self.node_types = {}
        self.next_id = 1

    def _encode_node(self, node: Dict[str, Any]) -> float:
        node_type = node.get("type", "Unknown")
        if node_type not in self.node_types:
            self.node_types[node_type] = self.next_id
            self.next_id += 1
        return float(self.node_types[node_type])

    def encode_tree(self, ast: Dict) -> np.ndarray:
        """Mã hóa AST bằng cách duyệt cây và gán chỉ số cho từng node."""
        encoded = []

        def traverse(n):
            encoded.append(self._encode_node(n))
            for child in n.get("children", []):
                traverse(child)

        traverse(ast)
        return np.array(encoded, dtype=np.float32)
