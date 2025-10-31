# =============================================================================
# OPTIMIZED MULTI-STAGE GPU-ENABLED DOCKERFILE FOR UNARYA AI MODEL SERVICE
# =============================================================================

# -----------------------------------------------------------------------------
# 1. Base Stage with Pre-installed PyTorch (SIGNIFICANTLY FASTER)
# -----------------------------------------------------------------------------
FROM pytorch/pytorch:2.9.0-cuda12.8-cudnn9-runtime AS pytorch-base

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# -----------------------------------------------------------------------------
# 2. Builder Stage (Install remaining dependencies)
# -----------------------------------------------------------------------------
FROM pytorch-base AS builder

WORKDIR /build

# Install build tools
RUN pip install --upgrade pip wheel setuptools

# Copy only requirements first (for layer caching)
COPY requirements.txt ./

# Create optimized requirements without torch (already installed in base image)
RUN grep -v "^torch" requirements.txt > requirements-optimized.txt || cp requirements.txt requirements-optimized.txt

# Install remaining dependencies and create wheels
RUN pip wheel --no-cache-dir -r requirements-optimized.txt -w /wheels

# -----------------------------------------------------------------------------
# 3. Development Stage (Fast rebuild for development)
# -----------------------------------------------------------------------------
FROM pytorch-base AS development

WORKDIR /app

ENV PYTHONUNBUFFERED=1

# Install development tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    git vim curl inotify-tools \
    && rm -rf /var/lib/apt/lists/*

# Install dependencies from wheels (MUCH FASTER than pip install)
COPY --from=builder /wheels /wheels
RUN pip install --no-cache-dir /wheels/*

COPY . .

# Clean build artifacts
RUN find . -type d -name "*.egg-info" -exec rm -rf {} + 2>/dev/null || true && \
    find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true && \
    find . -type d -name "*.dist-info" -exec rm -rf {} + 2>/dev/null || true && \
    find . -type d -name "build" -exec rm -rf {} + 2>/dev/null || true && \
    find . -type d -name ".eggs" -exec rm -rf {} + 2>/dev/null || true

# Install package in editable mode
RUN pip install --no-cache-dir -e .

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/bin/bash", "-c", "/app/entrypoint.sh"]

# -----------------------------------------------------------------------------
# 4. Production Stage (Minimal GPU runtime)
# -----------------------------------------------------------------------------
FROM nvidia/cuda:13.0.1-cudnn-runtime-ubuntu24.04 AS production

ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONUNBUFFERED=1
ENV PATH="/usr/local/bin:$PATH"
ENV MODEL_PORT=6000
ENV MODEL_NAME=unarya-intelligence

# Install Python and essential tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3.14 \
    python3-pip \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy pre-built wheels from builder
COPY --from=builder /wheels /wheels

# Install PyTorch with CUDA support (pre-built wheel is much faster)
RUN pip install --no-cache-dir \
    torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121

# Install remaining dependencies
RUN pip install --no-cache-dir /wheels/* && rm -rf /wheels

# Copy application code
COPY src/ ./src/
COPY setup.py ./
RUN pip install --no-cache-dir .

# Runtime volumes
VOLUME ["/models", "/data"]

# Security: non-root user
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 6000

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD python3 -c "import torch; print(torch.cuda.is_available())" || exit 1

ENTRYPOINT ["python3", "-u", "-m", "src.modelserver.server"]