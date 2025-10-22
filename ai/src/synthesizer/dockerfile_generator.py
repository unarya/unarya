from .schemas import CodeContext
from .template_engine import TemplateEngine

class DockerfileGenerator:
    """Sinh Dockerfile tối ưu cho từng ngôn ngữ."""

    def __init__(self):
        self.engine = TemplateEngine()

    def generate(self, context: CodeContext) -> str:
        base_images = {
            "python": "python:3.10-slim",
            "node": "node:18-alpine",
            "java": "eclipse-temurin:17-jdk",
            "go": "golang:1.21-alpine",
        }

        base = base_images.get(context.language.lower(), "ubuntu:22.04")

        template = """\
# Generated Dockerfile
FROM {{ base }} AS build
WORKDIR /app
COPY . .
{% if language == 'python' %}
RUN pip install -r requirements.txt
CMD ["python", "{{ entry_point }}"]
{% elif language == 'node' %}
RUN npm install
CMD ["node", "{{ entry_point }}"]
{% elif language == 'java' %}
RUN javac {{ entry_point }}
CMD ["java", "{{ entry_point.split('.')[0] }}"]
{% elif language == 'go' %}
RUN go build -o app {{ entry_point }}
CMD ["./app"]
{% else %}
CMD ["bash"]
{% endif %}
"""

        return self.engine.render(template, {
            "base": base,
            "language": context.language.lower(),
            "entry_point": context.entry_point
        })
