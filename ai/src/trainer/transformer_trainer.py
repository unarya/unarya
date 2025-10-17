import torch
import torch.nn as nn
import torch.optim as optim
from .types import Dataset, Config

class SimpleTransformer(nn.Module):
    """Một transformer rất đơn giản cho minh họa."""
    def __init__(self, vocab_size=1000, hidden_dim=128, num_layers=2):
        super().__init__()
        self.embed = nn.Embedding(vocab_size, hidden_dim)
        encoder_layer = nn.TransformerEncoderLayer(d_model=hidden_dim, nhead=4)
        self.encoder = nn.TransformerEncoder(encoder_layer, num_layers=num_layers)
        self.fc = nn.Linear(hidden_dim, vocab_size)

    def forward(self, x):
        embedded = self.embed(x)
        encoded = self.encoder(embedded)
        return self.fc(encoded)

class TransformerTrainer:
    """Huấn luyện mô hình hiểu mã (BERT-like hoặc GPT-like)."""

    def __init__(self):
        self.model = None

    def train(self, dataset: Dataset, config: Config):
        print(f"[TransformerTrainer] Training transformer with {len(dataset.examples)} samples")
        self.model = SimpleTransformer(vocab_size=5000, hidden_dim=config.embedding_dim)
        optimizer = optim.Adam(self.model.parameters(), lr=config.learning_rate)
        criterion = nn.CrossEntropyLoss()

        for epoch in range(config.epochs):
            total_loss = 0
            for ex in dataset.examples:
                x = torch.randint(0, 5000, (10, 1))  # dummy input
                y = torch.randint(0, 5000, (10, 1))
                optimizer.zero_grad()
                pred = self.model(x).view(-1, 5000)
                loss = criterion(pred, y.view(-1))
                loss.backward()
                optimizer.step()
                total_loss += loss.item()
            print(f"Epoch {epoch+1}/{config.epochs}, Loss: {total_loss:.4f}")

    def save_model(self, path: str):
        torch.save(self.model.state_dict(), path)
        print(f"[TransformerTrainer] Model saved to {path}")
