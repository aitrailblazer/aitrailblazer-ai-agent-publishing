#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${PUBLISHING_AGENT_URL:-http://localhost:8080}"
HEADER_ARGS=()
if [[ -n "${PUBLISHING_DEMO_API_KEY:-}" ]]; then
  HEADER_ARGS+=(-H "X-Demo-Key: ${PUBLISHING_DEMO_API_KEY}")
fi

curl -s "${BASE_URL}/v1/usage" "${HEADER_ARGS[@]}"
echo
