import time
from typing import Any, Dict
from .schemas import CacheEntry

class ResultCache:
    """Bộ nhớ đệm kết quả inference để tăng tốc phản hồi."""

    def __init__(self, ttl: int = 60):
        self.ttl = ttl
        self._store: Dict[str, CacheEntry] = {}

    def get(self, key: str) -> Any:
        entry = self._store.get(key)
        if entry and time.time() - entry.timestamp < self.ttl:
            return entry.value
        elif entry:
            del self._store[key]
        return None

    def set(self, key: str, value: Any):
        self._store[key] = CacheEntry(key=key, value=value, timestamp=time.time())

    def clear(self):
        self._store.clear()
