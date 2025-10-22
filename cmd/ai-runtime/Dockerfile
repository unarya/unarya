# ==========================================
# Dockerfile - AI Runtime with ONNX GPU
# ==========================================
FROM nvidia/cuda:13.0.1-cudnn-runtime-ubuntu24.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV GOLANG_VERSION=1.21.5
ENV ONNXRUNTIME_VERSION=1.16.3
ENV PATH=/usr/local/go/bin:$PATH
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

# Install system dependencies
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    git \
    build-essential \
    pkg-config \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Go
RUN wget https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

# Install ONNX Runtime with GPU support
RUN wget https://github.com/microsoft/onnxruntime/releases/download/v${ONNXRUNTIME_VERSION}/onnxruntime-linux-x64-gpu-${ONNXRUNTIME_VERSION}.tgz && \
    tar -xzf onnxruntime-linux-x64-gpu-${ONNXRUNTIME_VERSION}.tgz && \
    mv onnxruntime-linux-x64-gpu-${ONNXRUNTIME_VERSION} /opt/onnxruntime && \
    rm onnxruntime-linux-x64-gpu-${ONNXRUNTIME_VERSION}.tgz

# Copy ONNX Runtime libraries to system path
RUN cp /opt/onnxruntime/lib/* /usr/local/lib/ && \
    ldconfig

# Set ONNX Runtime environment
ENV ONNXRUNTIME_ROOT=/opt/onnxruntime
ENV CGO_CFLAGS="-I${ONNXRUNTIME_ROOT}/include"
ENV CGO_LDFLAGS="-L${ONNXRUNTIME_ROOT}/lib -lonnxruntime"

# Create app directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o ai_runtime ./cmd/main.go

# Create models directory
RUN mkdir -p /app/models

# Expose gRPC port
EXPOSE 6000

# Run the application
CMD ["./ai_runtime"]
