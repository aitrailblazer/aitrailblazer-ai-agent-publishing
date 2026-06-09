#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${PUBLISHING_AGENT_URL:-http://localhost:8080}"
AUTH_ARGS=()
if [[ -n "${PUBLISHING_DEMO_API_KEY:-}" ]]; then
  AUTH_ARGS+=(-H "X-Demo-Key: ${PUBLISHING_DEMO_API_KEY}")
fi
JSON_ARGS=("${AUTH_ARGS[@]}" -H "Content-Type: application/json")

echo "== health =="
curl -s "${BASE_URL}/health"
echo
echo

echo "== archive brief =="
curl -s "${BASE_URL}/v1/archive-brief" \
  "${JSON_ARGS[@]}" \
  -d '{"publication_url":"https://deltasignal.substack.com/","publication":"DeltaSignal","question":"How can this archive become agent-readable memory for other publishers?"}'
echo
echo

echo "== tripcode resolve =="
curl -s "${BASE_URL}/v1/tripcode" \
  "${JSON_ARGS[@]}" \
  -d '{"tripcode":"HUT-RIVER-001","session_id":"demo","include_river":true,"include_mongo_fit":true}'
echo
echo

echo "== judge proof package =="
curl -s "${BASE_URL}/v1/judge-demo" \
  "${AUTH_ARGS[@]}"
echo
echo

echo "== session follow-up =="
curl -s -X POST "${BASE_URL}/resolve?session_id=demo" \
  "${AUTH_ARGS[@]}" \
  -H "Content-Type: text/plain" \
  --data 'Using the previous publication River, what should stay in context?'
echo
