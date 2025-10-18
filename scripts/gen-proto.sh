#!/usr/bin/env bash
#
# gen-proto.sh — Generate Go and Python gRPC bindings from .proto files.
# Fully self-contained, works even on a fresh environment.
#

set -e  # Exit immediately on error

# === CONFIGURATION ===
PROTO_DIR="./lib/proto"
GO_OUT_BASE="./lib/proto/pb"
PY_OUT_DIR="./ai/src/pb"
PY_ENV_DIR="./ai/.venv"
GOPATH_BIN="$(go env GOPATH)/bin"

# === COLORS ===
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo ""
echo -e "${YELLOW}=== 🧩 Unarya Proto Generator ===${NC}"

# === 1. CHECK protoc ===
echo -e "${YELLOW}→ Checking protoc...${NC}"
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}❌ protoc not found.${NC}"
    echo "Installing on Ubuntu: sudo apt install -y protobuf-compiler"
    echo "On macOS: brew install protobuf"
    exit 1
else
    echo -e "${GREEN}✓ protoc found:${NC} $(protoc --version)"
fi

# === 2. CHECK Go plugins ===
echo -e "${YELLOW}→ Checking Go gRPC plugins...${NC}"

install_go_plugin() {
    local pkg=$1
    local name=$2
    if ! command -v "$name" &> /dev/null; then
        echo "Installing $name..."
        go install "$pkg@latest"
    fi
}

install_go_plugin "google.golang.org/protobuf/cmd/protoc-gen-go" "protoc-gen-go"
install_go_plugin "google.golang.org/grpc/cmd/protoc-gen-go-grpc" "protoc-gen-go-grpc"

export PATH="$PATH:$GOPATH_BIN"
echo -e "${GREEN}✓ Go plugins ready.${NC}"

# === 3. PYTHON ENVIRONMENT SETUP ===
if command -v python3 &> /dev/null; then
    echo -e "${YELLOW}→ Checking Python environment...${NC}"

    if ! dpkg -s python3-venv &> /dev/null || ! dpkg -s python3-pip &> /dev/null; then
        echo "Installing missing Python packages (python3-venv, python3-pip, setuptools)..."
        sudo apt-get update -y >/dev/null
        sudo apt-get install -y python3-venv python3-pip python3-setuptools >/dev/null
    fi

    if [ ! -d "$PY_ENV_DIR" ]; then
        echo "Creating virtual environment at $PY_ENV_DIR..."
        python3 -m venv "$PY_ENV_DIR"
    fi

    source "$PY_ENV_DIR/bin/activate"

    python3 -m ensurepip --upgrade >/dev/null 2>&1 || true

    if ! python3 -m grpc_tools.protoc --version &> /dev/null; then
        echo "Installing grpcio + grpcio-tools into venv..."
        pip install --quiet --upgrade pip
        pip install --quiet grpcio grpcio-tools
    fi

    echo -e "${GREEN}✓ Python environment ready.${NC}"
else
    echo -e "${RED}Python3 not found.${NC} Skipping Python generation."
fi

# === 4. CLEAN OUTPUT DIRECTORIES ===
echo -e "${YELLOW}→ Preparing output directories...${NC}"
rm -rf "$GO_OUT_BASE" "$PY_OUT_DIR"
mkdir -p "$GO_OUT_BASE" "$PY_OUT_DIR"

# === 5. GENERATE GO FILES ===
echo -e "${YELLOW}→ Generating Go gRPC bindings...${NC}"

for file in "$PROTO_DIR"/*.proto; do
    base=$(basename "$file" .proto)
    pkg_dir="${GO_OUT_BASE}/${base}pb"

    mkdir -p "$pkg_dir"

    echo "   ↳ $base.proto → $pkg_dir"

    protoc \
        --proto_path="$PROTO_DIR" \
        --go_out="$pkg_dir" \
        --go_opt=paths=source_relative \
        --go-grpc_out="$pkg_dir" \
        --go-grpc_opt=paths=source_relative \
        "$file"
done

echo -e "${GREEN}✓ Go proto files generated successfully.${NC}"

# === 6. GENERATE PYTHON FILES ===
if command -v python3 &> /dev/null; then
    echo -e "${YELLOW}→ Generating Python gRPC bindings...${NC}"

    find "$PROTO_DIR" -name "*.proto" | while read -r file; do
        rel_path="${file#$PROTO_DIR/}"
        echo "   ↳ $rel_path"

        out_dir="$PY_OUT_DIR/$(dirname "$rel_path")"
        mkdir -p "$out_dir"

        python3 -m grpc_tools.protoc \
            -I"$PROTO_DIR" \
            --python_out="$PY_OUT_DIR" \
            --grpc_python_out="$PY_OUT_DIR" \
            "$file"
    done

    # === FIX PYTHON IMPORTS ===
    echo -e "${YELLOW}→ Fixing Python imports for relative imports...${NC}"

    find "$PY_OUT_DIR" -name "*_pb2*.py" | while read -r py_file; do
        # Fix: import ai_pb2 → from . import ai_pb2
        sed -i 's/^import \([a-zA-Z0-9_]*\)_pb2 as \([a-zA-Z0-9_]*\)__pb2$/from . import \1_pb2 as \2__pb2/g' "$py_file"
    done

    # === CREATE __init__.py FILES ===
    echo -e "${YELLOW}→ Creating __init__.py files...${NC}"
    find "$PY_OUT_DIR" -type d | while read -r dir; do
        touch "$dir/__init__.py"
    done

    echo -e "${GREEN}✓ Python proto files generated successfully.${NC}"
    deactivate || true
fi

# === 7. DONE ===
echo ""
echo -e "${GREEN}✅ All proto files generated successfully!${NC}"
echo -e "→ Go output: ${YELLOW}${GO_OUT_BASE}${NC}"
echo -e "→ Python output: ${YELLOW}${PY_OUT_DIR}${NC}"
