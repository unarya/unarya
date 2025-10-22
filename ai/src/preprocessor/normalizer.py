import numpy as np

class Normalizer:
    """Chuẩn hóa dữ liệu (min-max normalization)."""

    def normalize(self, matrix: np.ndarray) -> np.ndarray:
        if matrix.size == 0:
            return matrix
        min_val = matrix.min()
        max_val = matrix.max()
        if min_val == max_val:
            return np.zeros_like(matrix)
        return (matrix - min_val) / (max_val - min_val)
