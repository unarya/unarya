#!/bin/bash
set -e

# Ensure the main binary is built before running
if [ ! -f "./bin/security_scan" ]; then
  echo "Binary './security_scan' not found. Building it now..."
  go build -buildvcs=false -o ./bin/security_scan .
fi

# Start the application with Air
echo "🚀 Starting the application with Air..."
exec air
