from .template_engine import TemplateEngine

class CIGenerator:
    """Sinh file cấu hình CI/CD cơ bản (GitHub Actions hoặc GitLab CI)."""

    def __init__(self):
        self.engine = TemplateEngine()

    def generate_github_actions(self, language: str) -> str:
        template = """\
name: Build and Test

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up {{ language }}
        uses: actions/setup-{{ language }}@v3
      - name: Install dependencies
        run: |
          {% if language == 'python' %}
          pip install -r requirements.txt
          {% elif language == 'node' %}
          npm install
          {% endif %}
      - name: Run tests
        run: |
          {% if language == 'python' %}
          pytest
          {% elif language == 'node' %}
          npm test
          {% endif %}
"""
        return self.engine.render(template, {"language": language.lower()})
