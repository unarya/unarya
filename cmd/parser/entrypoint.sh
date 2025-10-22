#!/bin/bash
set -e

# Ensure the main binary is built before running
if [ ! -f "./bin/parser" ]; then
  echo "Binary './parser' not found. Building it now..."
  go build -buildvcs=false -o ./bin/parser .
fi

# Start the application with Air
echo "🚀 Starting the application with Air..."
exec air
