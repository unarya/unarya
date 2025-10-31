#!/bin/bash
set -e
set -o pipefail

# ==============================================================================
# UNARYA - DEPLOY SCRIPT
# ------------------------------------------------------------------------------
# Description : Build & launch Unarya microservices (Collector, Orchestrator,
#               Parser, Security Scan, AI Model, Infra)
# Author      : Tiecont
# Version     : 1.0
# ==============================================================================

# ==== Configuration ====
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$ROOT_DIR/infra"
COLLECTOR_DIR="$ROOT_DIR/cmd/collector"
ORCHESTRATOR_DIR="$ROOT_DIR/cmd/orchestrator"
PARSER_DIR="$ROOT_DIR/cmd/parser"
SECURITY_SCAN_DIR="$ROOT_DIR/cmd/security_scan"
AI_DIR="$ROOT_DIR/ai"
ENV_DIR="$ROOT_DIR/configs/.env"
COLLECTOR_IMAGE_TAG="unarya-collector:dev"
ORCHESTRATOR_IMAGE_TAG="unarya-orchestrator:dev"
PARSER_IMAGE_TAG="unarya-parser:dev"
SECURITY_IMAGE_TAG="unarya-security_scan:dev"
AI_IMAGE_TAG="unarya-ai:dev"

# ==== Utilities ====
log() {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}

error_exit() {
    echo -e "\033[1;31m[ERROR]\033[0m $1"
    exit 1
}

# ==== Preflight ====
for dir in "$INFRA_DIR" "$COLLECTOR_DIR" "$ORCHESTRATOR_DIR" "$PARSER_DIR" "$SECURITY_SCAN_DIR" "$AI_DIR"; do
    [ -d "$dir" ] || error_exit "Missing directory: $dir"
done


if [ ! -f "$ENV_DIR" ]; then
  echo "Warning: $ENV_DIR not found. Using defaults from compose."
else
  echo "Found environment file: $ENV_DIR"
fi

# Load env vars for the execution
if [ -f "$ENV_DIR" ]; then
  source "$ENV_DIR" 2>/dev/null || true
fi

cd "$ROOT_DIR"

# ----------------------------------------------------------------------
# CLEANUP STAGE
# ----------------------------------------------------------------------
log "Stopping and cleaning up old containers..."
cd "$INFRA_DIR" && docker compose --env-file="$ENV_DIR" down --remove-orphans || error_exit "Failed to stop old containers."

# ----------------------------------------------------------------------
# BUILD STAGE
# ----------------------------------------------------------------------
build_service() {
    local name=$1
    local dir=$2
    local image=$3
    log "Building $name image from: $dir"
    cd "$dir"
    docker build \
        --target=development \
        -t "$image" \
        --build-arg CACHEBUST=$(date +%s) \
        . || error_exit "Failed to build $name image."
}

build_service "Collector" "$COLLECTOR_DIR" "$COLLECTOR_IMAGE_TAG"
build_service "Orchestrator" "$ORCHESTRATOR_DIR" "$ORCHESTRATOR_IMAGE_TAG"
build_service "Parser" "$PARSER_DIR" "$PARSER_IMAGE_TAG"
build_service "Security Scan" "$SECURITY_SCAN_DIR" "$SECURITY_IMAGE_TAG"
build_service "AI Model" "$AI_DIR" "$AI_IMAGE_TAG"

# ----------------------------------------------------------------------
# DEPLOY STAGE
# ----------------------------------------------------------------------
log "Starting infrastructure stack..."
cd "$INFRA_DIR"
docker compose --env-file "$ENV_DIR" up -d || error_exit "Failed to start infra stack."

log "Listing running containers..."
docker ps --filter "name=unarya-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "✅ Unarya microservices are running successfully!"
echo "------------------------------------------------"
echo "  • Collector Image:     $COLLECTOR_IMAGE_TAG"
echo "  • Orchestrator Image:  $ORCHESTRATOR_IMAGE_TAG"
echo "  • Parser Image:        $PARSER_IMAGE_TAG"
echo "  • Security Scan Image: $SECURITY_IMAGE_TAG"
echo "  • AI Model Image:      $AI_IMAGE_TAG"
echo "  • Infra Compose:       $INFRA_DIR"
echo "------------------------------------------------"
echo ""
