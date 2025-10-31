import yaml
from .schemas import ValidationResult

class Validator:
    """Kiểm tra tính hợp lệ của các file cấu hình sinh ra."""

    def validate_yaml(self, content: str) -> ValidationResult:
        try:
            yaml.safe_load(content)
            return ValidationResult(valid=True, errors=[])
        except yaml.YAMLError as e:
            return ValidationResult(valid=False, errors=[str(e)])

    def validate_dockerfile(self, content: str) -> ValidationResult:
        errors = []
        if not content.startswith("FROM"):
            errors.append("Dockerfile must start with a FROM instruction.")
        if "CMD" not in content:
            errors.append("Dockerfile missing CMD instruction.")
        return ValidationResult(valid=(len(errors) == 0), errors=errors)
