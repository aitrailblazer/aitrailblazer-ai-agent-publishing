#!/usr/bin/env bash
set -euo pipefail

missing=0

if [[ -z "${MDB_MCP_CONNECTION_STRING:-}" && ( -z "${MDB_MCP_API_CLIENT_ID:-}" || -z "${MDB_MCP_API_CLIENT_SECRET:-}" ) ]]; then
  echo "missing MongoDB MCP credentials: set MDB_MCP_CONNECTION_STRING or both MDB_MCP_API_CLIENT_ID and MDB_MCP_API_CLIENT_SECRET" >&2
  missing=1
fi

if [[ -z "${MCP_SERVER_URL:-}" ]]; then
  echo "missing MCP_SERVER_URL for the app runtime proof path" >&2
  missing=1
fi

if ! npx --yes mongodb-mcp-server@latest --help >/dev/null 2>&1; then
  echo "mongodb-mcp-server@latest is not runnable through npx" >&2
  missing=1
fi

if [[ "${missing}" -ne 0 ]]; then
  exit 1
fi

echo "MongoDB MCP prerequisites present."
