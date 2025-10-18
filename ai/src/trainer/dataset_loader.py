import json
from typing import List
from .schemas import Dataset, Example

class DatasetLoader:
    """Tải tập dữ liệu huấn luyện từ tệp hoặc thư mục."""

    def load_from_json(self, path: str) -> Dataset:
        with open(path, "r", encoding="utf-8") as f:
            data = json.load(f)
        examples = [Example(**item) for item in data]
        return Dataset(examples=examples)

    def load_from_texts(self, inputs: List[str], outputs: List[str]) -> Dataset:
        examples = [Example(input_text=i, output_text=o) for i, o in zip(inputs, outputs)]
        return Dataset(examples=examples)

    def sample(self, dataset: Dataset, n: int) -> Dataset:
        """Trích xuất một mẫu nhỏ để huấn luyện nhanh."""
        return Dataset(examples=dataset.examples[:n])
