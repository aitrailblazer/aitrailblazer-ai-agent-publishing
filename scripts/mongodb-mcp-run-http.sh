#!/usr/bin/env bash
set -euo pipefail

PORT="${MONGODB_MCP_PORT:-3000}"
HOST="${MONGODB_MCP_HOST:-127.0.0.1}"

if [[ -z "${MDB_MCP_CONNECTION_STRING:-}" && ( -z "${MDB_MCP_API_CLIENT_ID:-}" || -z "${MDB_MCP_API_CLIENT_SECRET:-}" ) ]]; then
  echo "Set MDB_MCP_CONNECTION_STRING or Atlas service-account env vars before starting MongoDB MCP." >&2
  exit 1
fi

export MDB_MCP_READ_ONLY="${MDB_MCP_READ_ONLY:-true}"
export MDB_MCP_TELEMETRY="${MDB_MCP_TELEMETRY:-disabled}"
export MDB_MCP_EXTERNALLY_MANAGED_SESSIONS="${MDB_MCP_EXTERNALLY_MANAGED_SESSIONS:-true}"
export MDB_MCP_HTTP_RESPONSE_TYPE="${MDB_MCP_HTTP_RESPONSE_TYPE:-json}"

exec npx --yes mongodb-mcp-server@latest \
	--transport http \
	--httpHost="${HOST}" \
	--httpPort="${PORT}" \
	--externallyManagedSessions \
	--httpResponseType=json \
	--readOnly
