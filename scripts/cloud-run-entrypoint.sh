#!/usr/bin/env bash
set -euo pipefail

APP_BIN="/aitrailblazer-ai-agent-publishing"
MONGO_PORT="${MONGODB_PORT:-27017}"
MCP_HOST="${MONGODB_MCP_HOST:-127.0.0.1}"
MCP_PORT="${MONGODB_MCP_PORT:-3000}"
DB_NAME="${MONGODB_DATABASE:-aitrailblazer_demo}"

export MONGODB_DATABASE="$DB_NAME"
export MDB_MCP_CONNECTION_STRING="${MDB_MCP_CONNECTION_STRING:-mongodb://127.0.0.1:${MONGO_PORT}/?directConnection=true}"
export MDB_MCP_READ_ONLY="${MDB_MCP_READ_ONLY:-true}"
export MDB_MCP_TELEMETRY="${MDB_MCP_TELEMETRY:-disabled}"
export MDB_MCP_EXTERNALLY_MANAGED_SESSIONS="${MDB_MCP_EXTERNALLY_MANAGED_SESSIONS:-true}"
export MDB_MCP_HTTP_RESPONSE_TYPE="${MDB_MCP_HTTP_RESPONSE_TYPE:-json}"
export MCP_SERVER_URL="${MCP_SERVER_URL:-http://127.0.0.1:${MCP_PORT}/mcp}"
export MCP_METHOD="${MCP_METHOD:-tools/call}"
export MCP_TOOL_NAME="${MCP_TOOL_NAME:-find}"
export MCP_SESSION_ID="${MCP_SESSION_ID:-aitrailblazer-judge-proof}"

children=()
cleanup() {
  for pid in "${children[@]:-}"; do
    kill "$pid" 2>/dev/null || true
  done
}
trap cleanup EXIT INT TERM

wait_for_mongo() {
  for _ in $(seq 1 60); do
    if mongosh "$MDB_MCP_CONNECTION_STRING" --quiet --eval 'db.runCommand({ ping: 1 }).ok' >/tmp/mongo-ping.out 2>/tmp/mongo-ping.err; then
      return 0
    fi
    sleep 1
  done
  echo "MongoDB did not become ready" >&2
  cat /tmp/mongo-ping.err >&2 || true
  return 1
}

wait_for_mcp() {
  for _ in $(seq 1 60); do
    if curl -fsS -m 2 \
      -H 'Content-Type: application/json' \
      -H 'Accept: application/json, text/event-stream' \
      -H "mcp-session-id: ${MCP_SESSION_ID}" \
      --data '{"jsonrpc":"2.0","id":"startup","method":"tools/list","params":{}}' \
      "$MCP_SERVER_URL" >/tmp/mcp-tools.out 2>/tmp/mcp-tools.err; then
      return 0
    fi
    sleep 1
  done
  echo "MongoDB MCP did not become ready" >&2
  cat /tmp/mcp-tools.err >&2 || true
  return 1
}

seed_mongo() {
  mongosh "$MDB_MCP_CONNECTION_STRING" --quiet --eval "
const database = db.getSiblingDB('${DB_NAME}');
database.tripcodes.updateOne(
  { tripcode: 'HUT-RIVER-001' },
  { \$set: {
    tripcode: 'HUT-RIVER-001',
    code: 'HUT-RIVER-001',
    target_type: 'river',
    target_id: 'hut8-rerate-river',
    publication_id: 'deltasignal',
    title: 'Hut 8 River memory path',
    summary: 'Seeded judge-demo TripCode used by the publishing-memory agent.',
    river_id: 'hut8-rerate-river',
    updated_at: new Date()
  }},
  { upsert: true }
);
database.articles.updateOne(
  { article_id: 'hut8-rerate-anchor' },
  { \$set: {
    article_id: 'hut8-rerate-anchor',
    tripcode: 'HUT-RIVER-001',
    title: 'Anchor TripCode article',
    publication_id: 'deltasignal',
    river_id: 'hut8-rerate-river',
    claims: ['prior thesis', 'anchor update', 'monitor next'],
    updated_at: new Date()
  }},
  { upsert: true }
);
database.river_edges.updateOne(
  { from: 'hut8-rerate-prior', to: 'hut8-rerate-anchor' },
  { \$set: {
    from: 'hut8-rerate-prior',
    to: 'hut8-rerate-anchor',
    relationship: 'continuation',
    reason: 'Judge-demo River continuity edge.',
    updated_at: new Date()
  }},
  { upsert: true }
);
database.reader_sessions.updateOne(
  { session_id: 'judge-demo' },
  { \$set: {
    session_id: 'judge-demo',
    last_tripcode: 'HUT-RIVER-001',
    monitor_next: ['new filings', 'hash price', 'capital allocation'],
    updated_at: new Date()
  }},
  { upsert: true }
);
database.agent_runs.updateOne(
  { run_id: 'judge-demo' },
  { \$set: {
    run_id: 'judge-demo',
    tool_calls: ['mongodb-mcp:find', 'gemini:generateContent', 'agent-builder:search'],
    model_lane: 'gemini-2.5-flash',
    summary: 'Seeded runtime proof record.',
    updated_at: new Date()
  }},
  { upsert: true }
);
database.claims.updateOne(
  { claim_id: 'hut8-rerate-claim-1' },
  { \$set: {
    claim_id: 'hut8-rerate-claim-1',
    article_id: 'hut8-rerate-anchor',
    claim_text: 'TripCodes let a reader resolve a publication memory path.',
    supporting_sources: ['hut8-rerate-anchor'],
    updated_at: new Date()
  }},
  { upsert: true }
);
database.tripcodes.createIndex({ tripcode: 1 }, { unique: true });
database.articles.createIndex({ tripcode: 1 });
database.reader_sessions.createIndex({ session_id: 1 }, { unique: true });
"
}

mkdir -p /tmp/mongodb
mongod --dbpath /tmp/mongodb --bind_ip 127.0.0.1 --port "$MONGO_PORT" --quiet --logpath /tmp/mongodb.log &
children+=("$!")
wait_for_mongo
seed_mongo

mongodb-mcp-server \
  --transport http \
  --httpHost "$MCP_HOST" \
  --httpPort "$MCP_PORT" \
  --externallyManagedSessions \
  --httpResponseType=json &
children+=("$!")
wait_for_mcp

"$APP_BIN" &
children+=("$!")

wait -n "${children[@]}"
