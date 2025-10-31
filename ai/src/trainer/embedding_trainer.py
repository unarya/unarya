import numpy as np
import os
from .schemas import Dataset, Config

class EmbeddingTrainer:
    """Huấn luyện mô hình embedding cho mã nguồn."""

    def __init__(self):
        self.embeddings = None

    def train(self, dataset: Dataset, config: Config):
        print(f"[EmbeddingTrainer] Training with {len(dataset.examples)} samples...")
        vocab = set()
        for ex in dataset.examples:
            vocab.update(ex.input_text.split())

        vocab = list(vocab)
        np.random.seed(42)
        self.embeddings = {
            token: np.random.randn(config.embedding_dim) for token in vocab
        }
        print(f"Generated embeddings for {len(vocab)} tokens with dim={config.embedding_dim}")

    def save_model(self, path: str):
        os.makedirs(os.path.dirname(path), exist_ok=True)
        np.savez(path, **self.embeddings)
        print(f"[EmbeddingTrainer] Model saved to {path}")
