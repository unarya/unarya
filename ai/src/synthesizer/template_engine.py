from jinja2 import Template

class TemplateEngine:
    """Công cụ render template Jinja2 cho các file cấu hình."""

    def render(self, template_str: str, context: dict) -> str:
        template = Template(template_str)
        return template.render(**context)
