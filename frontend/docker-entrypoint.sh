#!/bin/sh
set -eu

cat > /app/dist/config.js <<EOF
window.__HAPPYHOUSE_CONFIG__ = {
  apiBaseUrl: "${VITE_API_BASE_URL:-}"
};
EOF

exec sh -c "serve -s dist -l ${PORT:-3000}"
