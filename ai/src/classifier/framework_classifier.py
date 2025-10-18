import re
from typing import List
from .schemas import Framework

class FrameworkClassifier:
    """Nhận diện framework phổ biến dựa trên nội dung mã."""

    FRAMEWORK_PATTERNS = {
        "React": [r"from\s+['\"]react['\"]", r"ReactDOM\.render"],
        "Angular": [r"@Component", r"angular\.module"],
        "Vue": [r"new\s+Vue", r"v-bind:"],
        "Django": [r"from\s+django", r"urlpatterns"],
        "Flask": [r"from\s+flask", r"Flask\("],
        "Spring Boot": [r"@SpringBootApplication", r"org\.springframework"],
        "Express": [r"require\(['\"]express['\"]\)", r"app\.listen"],
    }

    def detect_frameworks(self, code: str) -> List[Framework]:
        detected = []
        for name, patterns in self.FRAMEWORK_PATTERNS.items():
            confidence = 0.0
            for pat in patterns:
                if re.search(pat, code):
                    confidence += 0.5
            if confidence > 0:
                detected.append(Framework(name=name, confidence=min(confidence, 1.0)))
        return detected
