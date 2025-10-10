# 🚀 Univia & Unarya

> **Intelligent Developer Hub powered by AI-driven Infrastructure Engine**

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![AI Powered](https://img.shields.io/badge/AI-Powered-blueviolet?logo=openai)](https://github.com)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

## 📖 Overview

**Univia & Unarya** is a next-generation developer platform that combines intelligent code analysis with automated infrastructure generation. The system consists of two core components:

- **🧩 Univia** - Developer Hub & Orchestration Platform
- **🧠 Unarya** - AI Engine for Code Understanding & Infrastructure Generation

### 🎯 Key Features

✅ **Intelligent Code Analysis** - Automatically detect languages, frameworks, and dependencies  
✅ **AI-Powered Dockerfile Generation** - Create optimized Docker configurations  
✅ **CI/CD Automation** - Generate pipeline configurations for GitHub Actions, GitLab CI  
✅ **Real-time Collaboration** - Built-in voice chat and meeting capabilities  
✅ **Team Management** - Advanced RBAC and organization controls  
✅ **GitHub Integration** - Seamless repo import and webhook synchronization  

---

## 🏗️ Architecture

```
                ┌──────────────────────────────┐
                │       Frontend (UI)          │
                │     React / Next.js          │
                └──────────────┬───────────────┘
                               │
                               ▼
                     ┌─────────────────┐
                     │   Univia API    │
                     │   (Go + Gin)    │
                     └────────┬────────┘
                              │
      ┌───────────────────────┼───────────────────────┐
      │                       │                       │
      ▼                       ▼                       ▼
┌───────────┐         ┌──────────────┐        ┌────────────────┐
│Auth/User  │         │Team Service  │        │Meeting Service │
│(JWT+Redis)│         │(RBAC, Org)   │        │(WebRTC Signal) │
└───────────┘         └──────────────┘        └────────────────┘
                              │
                              ▼
                    ┌────────────────┐
                    │ GitHub Import  │──→ Repo Metadata
                    │ Webhook Sync   │──→ Commit History
                    └────────────────┘
                              │
                              ▼
                  ┌────────────────────┐
                  │ CI/CD Orchestrator │──→ Build Queue
                  │  (Univia Runner)   │──→ Docker Build
                  └────────┬───────────┘
                           │
                           ▼
                  ┌────────────────────┐
                  │   Unarya Engine    │
                  │   (AI + Parser)    │
                  └────────────────────┘
                           │
                           ▼
                 ┌───────────────────────┐
                 │   uryad Core Engine   │
                 │  - Code Classifier    │
                 │  - Docker Synthesizer │
                 │  - CMD Predictor      │
                 └───────────────────────┘
```

---

## 🧩 Component Overview

### 1. **Univia** (Developer Hub)

**Purpose:** Manage user lifecycle, teams, projects, and CI/CD integration.

**Core Features:**

| Feature | Description |
|---------|-------------|
| 👤 **User Authentication** | Email, OAuth (GitHub, Google, etc.) |
| 👥 **Team Management** | Create teams, role-based permissions, repo sharing |
| 🐙 **GitHub Integration** | Import repositories, webhook events |
| 🧠 **Project Lifecycle** | Build pipelines, CI/CD triggers, version history |
| 📞 **Real-time Collaboration** | Voice chat, meetings, WebRTC signaling |
| 🚀 **Deployment Integration** | Trigger Docker builds, push, and deploy |

**Result:** Univia serves as the **frontend orchestration platform** where all user actions flow through and synchronize with Unarya.

---

### 2. **Unarya** (AI Engine Hub)

**Purpose:** Provide intelligence to Univia - "understand code, generate Dockerfiles, predict runtime environments."

**Core Features:**

| Feature | Description |
|---------|-------------|
| 🧬 **Source Code Understanding** | Classify projects (Go, Node, Python, Java, React, etc.) |
| 📦 **Dependency Extraction** | Parse package.json, go.mod, requirements.txt |
| 🧠 **AI Model Engine (uryad)** | Generate Dockerfile, .dockerignore, docker-compose.yaml |
| ⚙️ **Engine Registry** | Language-specific engines (GoEngine, NodeEngine, etc.) |
| 🔄 **Runtime Adapter** | Bi-directional communication with Univia via gRPC/WebSocket |

**Result:** Unarya is the **brain** - it reads, understands code structure, and generates runtime environments and build logic for Univia.

---

## ⚙️ Unarya Engine Components

| # | Engine | Role | Technology Stack |
|---|--------|------|------------------|
| 1️⃣ | **Code Classification** | Identify language & framework from source | Tree-sitter AST + Fine-tuned Transformer |
| 2️⃣ | **Dependency Extraction** | Analyze config files (go.mod, package.json, etc.) | Rule-based parser + Static analysis |
| 3️⃣ | **Docker Synthesizer** | Generate optimized Dockerfiles (multi-stage builds) | CodeLlama, Phi-3, StarCoder2 |
| 4️⃣ | **Environment Predictor** | Predict CMD, ENV variables needed to run app | Language-specific heuristics + Embeddings |
| 5️⃣ | **CI/CD Generator** | Generate `.github/workflows` or `.gitlab-ci.yml` | Template engine + LLM-guided |
| 6️⃣ | **Security Scanner** | Detect vulnerabilities, outdated packages | Snyk-like engine + Custom vulnerability DB |
| 7️⃣ | **uryad Runtime Adapter** | Bridge between Univia and engines | gRPC microservice + Redis cache + Kafka queue |

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21+
- **Docker** & Docker Compose
- **Node.js** 18+ (for frontend)
- **Redis** (for caching)
- **PostgreSQL** (for data storage)

### Installation

```bash
# Clone the repository
git clone https://github.com/your-org/univia-unarya.git
cd univia-unarya

# Start services with Docker Compose
docker-compose up -d

# Install Univia dependencies
cd univia
go mod download

# Install Unarya dependencies
cd ../unarya
go mod download

# Run Univia API
cd ../univia
go run cmd/api/main.go

# Run Unarya Engine
cd ../unarya
go run cmd/engine/main.go
```

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

Access the application at `http://localhost:3000`

---

## 🛠️ Development

### Project Structure

```
.
├── univia/                 # Developer Hub (Orchestration)
│   ├── cmd/               # Entry points
│   ├── internal/          # Core business logic
│   │   ├── auth/         # Authentication service
│   │   ├── team/         # Team management
│   │   ├── github/       # GitHub integration
│   │   ├── cicd/         # CI/CD orchestrator
│   │   └── meeting/      # WebRTC signaling
│   └── pkg/              # Shared packages
│
├── unarya/                # AI Engine Hub
│   ├── cmd/              # Entry points
│   ├── internal/         # Core AI engines
│   │   ├── classifier/   # Code classification
│   │   ├── parser/       # Dependency extraction
│   │   ├── synthesizer/  # Docker generation
│   │   ├── predictor/    # Environment prediction
│   │   └── scanner/      # Security scanning
│   └── pkg/              # Shared packages
│
├── frontend/             # React/Next.js UI
├── docker-compose.yml    # Development environment
└── README.md            # This file
```

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards

- Follow Go best practices and `gofmt` formatting
- Write unit tests for new features
- Update documentation for API changes
- Use conventional commits

---

## 🌟 Roadmap

- [x] Core Univia API
- [x] GitHub Integration
- [x] Basic AI Engine (Code Classification)
- [ ] Advanced Docker Synthesizer with LLM
- [ ] Security Vulnerability Scanner
- [ ] Kubernetes Deployment Support
- [ ] Multi-cloud Provider Integration
- [ ] Plugin System for Custom Engines

---

## 📞 Contact & Support

- **Documentation:** [docs.univia.dev](https://docs.univia.dev)
- **Issues:** [GitHub Issues](https://github.com/your-org/univia-unarya/issues)
- **Discussions:** [GitHub Discussions](https://github.com/your-org/univia-unarya/discussions)
- **Email:** support@univia.dev

---

<div align="center">

**Made with ❤️ by the Univia & Unarya Team**

⭐ Star us on GitHub — it helps!

</div>
