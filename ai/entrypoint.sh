#!/usr/bin/env bash
set -e

echo "🚀 Starting Unarya AI Model Server (Dev Mode, GPU)"
echo "📁 Using hot reload script: /app/hot_reload.py"

# Đảm bảo thư mục code tồn tại
if [ ! -d "/app/src" ]; then
  echo "❌ Error: /app/src not found. Check your docker volume mapping."
  exit 1
fi

# Cài watchdog nếu hot_reload cần mà chưa có
if ! python3 -c "import watchdog" &> /dev/null; then
  echo "📦 Installing watchdog..."
  pip install --no-cache-dir watchdog >/dev/null
fi

# Chạy hot reload
exec python3 -u /app/hot_reload.py
