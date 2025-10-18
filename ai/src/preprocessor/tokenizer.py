import re
from typing import List
from .schemas import Token

class CodeTokenizer:
    """Tokenizer cơ bản cho mã nguồn Python."""

    def __init__(self):
        # Regex đơn giản cho định danh, số, ký tự đặc biệt
        self.pattern = re.compile(r"[A-Za-z_]\w*|\d+|==|!=|<=|>=|[^\s]")

    def tokenize(self, code: str) -> List[Token]:
        tokens = []
        for pos, match in enumerate(self.pattern.finditer(code)):
            value = match.group()
            if value.isidentifier():
                tok_type = "IDENTIFIER"
            elif value.isdigit():
                tok_type = "NUMBER"
            elif value in {"=", "==", "+", "-", "*", "/", "(", ")", "{", "}", "[", "]", ":", ","}:
                tok_type = "SYMBOL"
            else:
                tok_type = "UNKNOWN"
            tokens.append(Token(type=tok_type, value=value, position=pos, metadata={}))
        return tokens
