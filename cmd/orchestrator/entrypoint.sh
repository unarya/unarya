#!/bin/bash
set -e

# Ensure the main binary is built before running
if [ ! -f "./bin/orchestrator" ]; then
  echo "Binary './orchestrator' not found. Building it now..."
  go build -buildvcs=false -o ./bin/orchestrator .
fi

# Start the application with Air
echo "🚀 Starting the application with Air..."
exec air
