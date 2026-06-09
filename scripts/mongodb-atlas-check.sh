#!/usr/bin/env bash
set -euo pipefail

missing=0

if [[ "${MONGODB_DEPLOYMENT_KIND:-}" != "atlas" ]]; then
  echo "MONGODB_DEPLOYMENT_KIND must be atlas for an Atlas runtime check." >&2
  missing=1
fi

case "${MDB_MCP_CONNECTION_STRING:-}" in
  mongodb+srv://*|mongodb://*)
    ;;
  *)
    echo "MDB_MCP_CONNECTION_STRING must be a MongoDB or MongoDB Atlas connection string." >&2
    missing=1
    ;;
esac

case "${MDB_MCP_CONNECTION_STRING:-}" in
  mongodb://127.0.0.1:*|mongodb://localhost:*|mongodb://0.0.0.0:*)
    echo "Atlas mode cannot use a localhost MongoDB connection string." >&2
    missing=1
    ;;
esac

if [[ -z "${MCP_SERVER_URL:-}" ]]; then
  echo "MCP_SERVER_URL is required for the app runtime proof path." >&2
  missing=1
fi

if [[ "${missing}" -ne 0 ]]; then
  exit 1
fi

echo "MongoDB Atlas runtime prerequisites present."
