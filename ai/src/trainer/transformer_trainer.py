import torch
import torch.nn as nn
import torch.optim as optim
from .schemas import Dataset, Config

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
    """Huấn luyện và export Transformer model (PyTorch → ONNX)."""

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

    def save_model(self, path: str, export_type: str = "onnx"):
        """
        Lưu model với định dạng mong muốn:
        - 'onnx'  → Dành cho runtime Golang (khuyên dùng)
        - 'torchscript-trace' → Dành cho inference Python/C++
        - 'torchscript-script' → Dành cho mô hình có control flow phức tạp
        - 'state_dict' → Dành cho nghiên cứu / fine-tune tiếp
        """
        self.model.eval()

        # Dummy input cho export
        dummy_input = torch.randint(0, 5000, (10, 1))

        if export_type == "onnx":
            export_path = f"{path}.onnx"
            torch.onnx.export(
                self.model,
                dummy_input,
                export_path,
                input_names=["input_ids"],
                output_names=["logits"],
                opset_version=17,
                dynamic_axes={"input_ids": {0: "batch"}, "logits": {0: "batch"}},
            )
            print(f"[TransformerTrainer] ✅ Model exported to ONNX: {export_path}")

        elif export_type == "torchscript-trace":
            export_path = f"{path}_trace.pt"
            traced = torch.jit.trace(self.model, dummy_input)
            traced.save(export_path)
            print(f"[TransformerTrainer] ✅ TorchScript (trace) saved: {export_path}")

        elif export_type == "torchscript-script":
            export_path = f"{path}_script.pt"
            scripted = torch.jit.script(self.model)
            scripted.save(export_path)
            print(f"[TransformerTrainer] ✅ TorchScript (script) saved: {export_path}")

        elif export_type == "state_dict":
            export_path = f"{path}.pth"
            torch.save(self.model.state_dict(), export_path)
            print(f"[TransformerTrainer] ✅ State dict saved: {export_path}")

        else:
            raise ValueError(f"Unknown export_type: {export_type}")

