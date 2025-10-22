from typing import List
import yaml
from .schemas import Service

class ComposeGenerator:
    """Sinh file docker-compose.yml."""

    def generate(self, services: List[Service]) -> str:
        compose_dict = {
            "version": "3.9",
            "services": {},
        }
        for svc in services:
            compose_dict["services"][svc.name] = {
                "image": svc.image,
                "ports": svc.ports,
            }
            if svc.environment:
                compose_dict["services"][svc.name]["environment"] = svc.environment
            if svc.volumes:
                compose_dict["services"][svc.name]["volumes"] = svc.volumes
            if svc.depends_on:
                compose_dict["services"][svc.name]["depends_on"] = svc.depends_on

        return yaml.dump(compose_dict, sort_keys=False)
