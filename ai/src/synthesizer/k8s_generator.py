import yaml
from .schemas import AppConfig

class K8sGenerator:
    """Sinh manifest Kubernetes: Deployment, Service, Ingress."""

    def generate_deployment(self, config: AppConfig) -> str:
        deployment = {
            "apiVersion": "apps/v1",
            "kind": "Deployment",
            "metadata": {"name": config.name},
            "spec": {
                "replicas": config.replicas,
                "selector": {"matchLabels": {"app": config.name}},
                "template": {
                    "metadata": {"labels": {"app": config.name}},
                    "spec": {
                        "containers": [{
                            "name": config.name,
                            "image": config.image,
                            "ports": [{"containerPort": p} for p in (config.ports or [80])],
                            "env": [{"name": k, "value": v} for k, v in (config.env or {}).items()]
                        }]
                    }
                }
            }
        }
        return yaml.dump(deployment, sort_keys=False)

    def generate_service(self, config: AppConfig) -> str:
        service = {
            "apiVersion": "v1",
            "kind": "Service",
            "metadata": {"name": f"{config.name}-svc"},
            "spec": {
                "selector": {"app": config.name},
                "ports": [{"port": p, "targetPort": p} for p in (config.ports or [80])]
            }
        }
        return yaml.dump(service, sort_keys=False)

    def generate_ingress(self, config: AppConfig) -> str:
        if not config.domain:
            return "# No domain specified — skipping ingress generation.\n"
        ingress = {
            "apiVersion": "networking.k8s.io/v1",
            "kind": "Ingress",
            "metadata": {"name": f"{config.name}-ingress"},
            "spec": {
                "rules": [{
                    "host": config.domain,
                    "http": {
                        "paths": [{
                            "path": "/",
                            "pathType": "Prefix",
                            "backend": {
                                "service": {
                                    "name": f"{config.name}-svc",
                                    "port": {"number": (config.ports or [80])[0]}
                                }
                            }
                        }]
                    }
                }]
            }
        }
        return yaml.dump(ingress, sort_keys=False)
