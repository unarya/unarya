# Unarya

An intelligent code analysis and deployment automation platform that combines static analysis, security scanning, and AI-powered infrastructure generation.

## Overview

Unarya is a dual-service architecture that processes source code through a comprehensive pipeline:
- **Golang Service**: Handles code collection, parsing, security scanning, and orchestration
- **Python Service**: Performs AI-based preprocessing, classification, model training, and synthesis of deployment artifacts

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Golang Service                          │
│  ┌──────────┐  ┌─────────┐  ┌─────────────────┐           │
│  │Collector │→ │ Parser  │→ │ Security Scan   │→          │
│  └──────────┘  └─────────┘  └─────────────────┘           │
│         ↓                                                   │
│  ┌──────────────────────────────────────────────┐          │
│  │           Orchestrator (gRPC)                │          │
│  └──────────────────────────────────────────────┘          │
└────────────────────────┬────────────────────────────────────┘
                         │ gRPC Protocol Buffers
                         ↓
┌─────────────────────────────────────────────────────────────┐
│                     Python Service                          │
│  ┌──────────────┐  ┌────────────┐  ┌──────────┐           │
│  │Preprocessor  │→ │Classifier  │→ │ Trainer  │→          │
│  └──────────────┘  └────────────┘  └──────────┘           │
│                                           ↓                 │
│  ┌──────────────┐  ┌────────────────────────────┐          │
│  │Model Server  │← │     Synthesizer            │          │
│  └──────────────┘  └────────────────────────────┘          │
│         (Generates: Dockerfile, Compose, K8s YAML)          │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
unarya/
├── cmd/                          # Golang service entrypoints
│   ├── orchestrator/             # Main pipeline orchestrator
│   ├── collector/                # Source code collection service
│   ├── parser/                   # Code parsing & AST analysis
│   └── security_scan/            # Security vulnerability scanner
│
├── internal/                     # Golang internal packages
│   ├── collector/                # Repository & source fetching logic
│   ├── parser/                   # Code structure analysis
│   ├── security_scan/            # Static analysis & vulnerability checks
│   ├── orchestrator/             # gRPC request coordination
│   └── shared/                   # Common utilities
│       ├── logging/              # Structured logging
│       ├── config/               # Configuration management
│       ├── grpc/                 # gRPC client stubs
│       └── auth/                 # Authentication utilities
│
├── ai/                           # Python service
│   ├── src/
│   │   ├── preprocessor/         # Data transformation & tokenization
│   │   ├── classifier/           # Language & framework detection
│   │   ├── trainer/              # Model training pipeline
│   │   ├── synthesizer/          # Deployment artifact generation
│   │   └── modelserver/          # Inference API server
│   ├── proto/                    # Python gRPC stubs
│   ├── requirements.txt
│   └── Dockerfile
│
├── lib/                          # Shared protocol definitions
│   └── proto/                    # gRPC protocol buffers
│       └── pipeline.proto        # Service contract definitions
│
├── configs/                      # Configuration files
│   ├── golang-service.yaml
│   ├── python-service.yaml
│   └── models.yaml
│
├── infra/                        # Infrastructure & deployment
│   ├── docker-compose.yml        # Local development setup
│   └── k8s/                      # Kubernetes manifests
│       ├── golang-deployment.yaml
│       ├── python-deployment.yaml
│       ├── grpc-gateway.yaml
│       └── secrets.yaml
│
├── scripts/                      # Automation scripts
│   ├── gen-proto.sh              # Protocol buffer compilation
│   ├── run-local.sh              # Local development launcher
│   └── benchmark.sh              # Performance testing
│
├── docs/                         # Documentation
│   ├── architecture.md           # System design & data flow
│   ├── model-structure.md        # AI model specifications
│   └── grpc-contract.md          # API definitions
│
├── Makefile                      # Build automation
├── LICENSE
└── README.md
```

## Features

### Golang Service
- **Code Collection**: Fetch source code from repositories, archives, or URLs
- **AST Parsing**: Analyze code structure, dependencies, and language detection
- **Security Scanning**: Static analysis, secret detection, permission checks, dependency vulnerabilities
- **Orchestration**: Coordinate pipeline stages via gRPC

### Python Service
- **Preprocessing**: Transform code into numerical matrices, tokenization, tree encoding
- **Classification**: Detect programming languages, frameworks, and coding styles
- **Model Training**: Train embeddings, transformers, and generative models
- **Synthesis**: Generate deployment artifacts (Dockerfiles, Docker Compose, Kubernetes YAML)
- **Model Server**: Inference API for real-time generation

## Prerequisites

- **Go** 1.21+
- **Python** 3.11+
- **Protocol Buffers** compiler (protoc)
- **Docker** & **Docker Compose** (for containerized deployment)
- **Make** (optional, for build automation)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/unarya.git
cd unarya
```

### 2. Generate Protocol Buffers

```bash
./scripts/gen-proto.sh
```

This compiles `.proto` files for both Go and Python services.

### 3. Install Dependencies

**Golang Service:**
```bash
cd cmd
go mod download
```

**Python Service:**
```bash
cd ai
pip install -r requirements.txt
```

### 4. Run Locally

Using the provided script:
```bash
./scripts/run-local.sh
```

Or using Docker Compose:
```bash
docker-compose -f infra/docker-compose.yml up --build
```

### 5. Using Make

```bash
# Build all services
make build

# Run tests
make test

# Generate proto files
make proto

# Run locally
make run-local
```

## Configuration

Configuration files are located in the `configs/` directory:

- `golang-service.yaml`: Golang service settings (ports, timeouts, security rules)
- `python-service.yaml`: Python service settings (model paths, inference parameters)
- `models.yaml`: AI model configurations

## API Usage

### Submit Code Analysis Request

```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "repo_url": "https://github.com/user/repo",
    "branch": "main",
    "output_format": "docker-compose"
  }'
```

### Get Analysis Results

```bash
curl http://localhost:8080/api/v1/results/{job_id}
```

## Development

### Adding a New Analysis Module

1. Define the service in `lib/proto/pipeline.proto`
2. Regenerate stubs: `./scripts/gen-proto.sh`
3. Implement in `internal/` (Go) or `ai/src/` (Python)
4. Update orchestrator to include the new stage

### Running Tests

```bash
# Golang tests
make test-go

# Python tests
make test-python

# Integration tests
make test-integration
```

### Benchmarking

```bash
./scripts/benchmark.sh
```

Tests throughput and latency under various loads.

## Deployment

### Docker Compose (Development)

```bash
docker-compose -f infra/docker-compose.yml up
```

### Kubernetes (Production)

```bash
kubectl apply -f infra/k8s/
```

This deploys:
- Golang service pods
- Python service pods
- gRPC gateway
- Required secrets and config maps

## Documentation

- [Architecture Overview](docs/architecture.md) - System design and data flow
- [Model Structure](docs/model-structure.md) - AI model input/output specifications
- [gRPC Contract](docs/grpc-contract.md) - Detailed API definitions

## Performance

- **Throughput**: ~1000 requests/minute (varies by code complexity)
- **Latency**:
    - Code parsing: ~200ms
    - Security scan: ~500ms
    - AI inference: ~2-5s
    - End-to-end: ~5-10s

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## Support

For issues, questions, or contributions, please open an issue on GitHub or contact the maintainers.

---

**Built with ❤️ using Go, Python, gRPC, and AI**