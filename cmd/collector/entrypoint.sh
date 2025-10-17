#!/bin/bash
set -e

# Ensure the main binary is built before running
if [ ! -f "./bin/collector" ]; then
  echo "Binary './collector' not found. Building it now..."
  go build -buildvcs=false -o ./bin/collector .
fi

# Start the application with Air
echo "🚀 Starting the application with Air..."
exec air
