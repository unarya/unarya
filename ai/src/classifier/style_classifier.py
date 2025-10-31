import re
from typing import List
from .schemas import StyleMetrics

class StyleClassifier:
    """Phân tích phong cách viết mã."""

    def analyze_style(self, code: str) -> StyleMetrics:
        lines = code.splitlines()
        non_empty = [l for l in lines if l.strip()]
        if not non_empty:
            return StyleMetrics(4, "unknown", 0, 0, 0)

        # Xác định số khoảng trắng đầu dòng trung bình
        indentations = [len(re.match(r"^\s*", l).group()) for l in non_empty]
        indentation_spaces = int(sum(indentations) / len(indentations))

        # Xác định quy ước đặt tên
        if re.search(r"[a-z]+_[a-z]+", code):
            naming = "snake_case"
        elif re.search(r"[a-z]+[A-Z][a-z]+", code):
            naming = "camelCase"
        elif re.search(r"class\s+[A-Z][a-zA-Z]+", code):
            naming = "PascalCase"
        else:
            naming = "unknown"

        # Các chỉ số khác
        avg_len = sum(len(l) for l in non_empty) / len(non_empty)
        comments = [l for l in lines if l.strip().startswith(("#", "//", "/*"))]
        comment_ratio = len(comments) / max(1, len(lines))

        # Ước lượng độ phức tạp dựa trên số lần xuất hiện từ khóa điều kiện
        complexity = len(re.findall(r"\b(if|for|while|switch|case|try)\b", code)) / max(1, len(lines))

        return StyleMetrics(
            indentation_spaces=indentation_spaces,
            naming_convention=naming,
            average_line_length=round(avg_len, 2),
            comment_ratio=round(comment_ratio, 3),
            complexity_estimate=round(complexity, 3)
        )
