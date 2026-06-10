#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${PUBLISHING_AGENT_URL:-https://aitrailblazer-ai-agent-publishing-rmycwek6ba-uc.a.run.app}"
PUBLIC_REPO_URL="${PUBLIC_REPO_URL:-https://github.com/aitrailblazer/aitrailblazer-ai-agent-publishing}"
PUBLIC_RAW_BASE="${PUBLIC_RAW_BASE:-https://raw.githubusercontent.com/aitrailblazer/aitrailblazer-ai-agent-publishing/main}"
DOC_PATH="docs/AITrailblazer_AI_Agent_Publishing_Public_Docs_2026_06_10.html"
SESSION_ID="submission-readiness-$(date +%s)"

tmpdir="$(mktemp -d)"
cleanup() { rm -rf "${tmpdir}"; }
trap cleanup EXIT

pass() { printf 'ok - %s\n' "$1"; }
fail() { printf 'FAIL - %s\n' "$1" >&2; exit 1; }
require_cmd() { command -v "$1" >/dev/null 2>&1 || fail "missing required command: $1"; }

http_get() {
  local url="$1"
  local out="$2"
  local headers="$3"
  curl -fsSL -D "${headers}" -o "${out}" "${url}"
}

assert_json_value() {
  local file="$1"
  local query="$2"
  local label="$3"
  jq -e "${query}" "${file}" >/dev/null || fail "${label}"
  pass "${label}"
}

require_cmd curl
require_cmd jq
require_cmd rg
require_cmd npx

echo "== local gates =="
make validate >/dev/null
pass "make validate"
make check >/dev/null
pass "make check"

echo
echo "== official MongoDB MCP prerequisite =="
MDB_MCP_CONNECTION_STRING="${MDB_MCP_CONNECTION_STRING:-mongodb://127.0.0.1:27017/?directConnection=true}" \
MCP_SERVER_URL="${MCP_SERVER_URL:-http://127.0.0.1:3000/mcp}" \
  make mcp-check >/dev/null
pass "make mcp-check with embedded demo connection"

if [[ "${SUBMISSION_READINESS_REQUIRE_ATLAS:-false}" == "true" || "${MONGODB_DEPLOYMENT_KIND:-}" == "atlas" ]]; then
  echo
  echo "== Atlas prerequisite =="
  make atlas-check >/dev/null
  pass "make atlas-check"
else
  echo
  echo "skip - Atlas prerequisite not required for this run"
fi

echo
echo "== public docs and local link hygiene =="
doc_count="$(find docs -maxdepth 1 -type f -name '*.html' | wc -l | tr -d ' ')"
[[ "${doc_count}" == "1" ]] || fail "docs folder must contain exactly one HTML file"
pass "docs folder contains one HTML file"

[[ -f "${DOC_PATH}" ]] || fail "canonical docs file is missing"
pass "canonical docs file exists"

if rg -n "href=[\"']\\./docs/" index.html README.html CHANGELOG.html | rg -v "${DOC_PATH}" >/dev/null; then
  rg -n "href=[\"']\\./docs/" index.html README.html CHANGELOG.html | rg -v "${DOC_PATH}" >&2
  fail "public HTML links to stale docs files"
fi
pass "index/README/changelog docs links target consolidated public docs"

echo
echo "== public secret and evidence scan =="
private_evidence_pattern="Co""dex|Open""AI|Chat""GPT|Clau""de|Anth""ropic|Cur""sor|Co""pilot|csrf""-token|Rapid_Agent_Hackathon_Rules_Source_2026_06_08""\\.raw\\.html|860""1137|AI""za|s""k-[A-Za-z0-9]|mong""odb\\+srv://[^ <)\"']+:[^ <)\"']+@"
if rg -n "${private_evidence_pattern}" \
  index.html README.html CHANGELOG.html Makefile Dockerfile cmd internal scripts docs .env.example >"${tmpdir}/scan.txt"; then
  cat "${tmpdir}/scan.txt" >&2
  fail "public scan found prohibited private evidence or obvious token pattern"
fi
pass "public scan found no prohibited private evidence terms or obvious token patterns"

echo
echo "== hosted HTTP smoke =="
http_get "${BASE_URL}/" "${tmpdir}/root.html" "${tmpdir}/root.headers"
rg -q "AITrailblazer AI Agent Publishing" "${tmpdir}/root.html" || fail "hosted root missing project title"
rg -q "Try the live proof" "${tmpdir}/root.html" || fail "hosted root missing live proof CTA"
rg -q "Public docs" "${tmpdir}/root.html" || fail "hosted root missing consolidated docs link"
rg -q "A short code printed in an article" "${tmpdir}/root.html" || fail "hosted root missing plain-language TripCode definition"
rg -q "Resolve handle" "${tmpdir}/root.html" || fail "hosted root missing agent action trace"
pass "hosted root"

http_get "${BASE_URL}/health" "${tmpdir}/health.json" "${tmpdir}/health.headers"
assert_json_value "${tmpdir}/health.json" '.ok == true' "hosted health"

http_get "${BASE_URL}/${DOC_PATH}" "${tmpdir}/docs.html" "${tmpdir}/docs.headers"
rg -q "Consolidated Public Documentation" "${tmpdir}/docs.html" || fail "hosted docs missing consolidated docs marker"
rg -q "Embedded XML Contract" "${tmpdir}/docs.html" || fail "hosted docs missing embedded XML marker"
pass "hosted consolidated docs"

http_get "${BASE_URL}/LICENSE" "${tmpdir}/license.txt" "${tmpdir}/license.headers"
rg -q "MIT License" "${tmpdir}/license.txt" || fail "hosted license missing MIT License"
pass "hosted license"

http_get "${BASE_URL}/v1/judge-demo" "${tmpdir}/judge.json" "${tmpdir}/judge.headers"
assert_json_value "${tmpdir}/judge.json" '.runtime_proof[] | select(.system == "Gemini" and .status == "live invoked")' "Gemini live proof"
assert_json_value "${tmpdir}/judge.json" '.runtime_proof[] | select(.system == "Agent Builder" and .status == "live invoked")' "Agent Builder live proof"
assert_json_value "${tmpdir}/judge.json" '.runtime_proof[] | select(.system == "MongoDB" and (.status | contains("live")))' "MongoDB live proof"
assert_json_value "${tmpdir}/judge.json" '.runtime_proof[] | select(.system == "MongoDB MCP" and .status == "live invoked")' "MongoDB MCP live proof"
assert_json_value "${tmpdir}/judge.json" '.tripcode_result.tripcode == "HUT-RIVER-001"' "judge demo TripCode result"

resolve_url="${BASE_URL}/resolve?tripcode=HUT-RIVER-001&session_id=${SESSION_ID}&question=What%20changed%20across%20this%20River%3F&include_river=true&include_mongo_fit=true"
http_get "${resolve_url}" "${tmpdir}/resolve.json" "${tmpdir}/resolve.headers"
assert_json_value "${tmpdir}/resolve.json" '.tripcode == "HUT-RIVER-001"' "resolve HUT-RIVER-001"
assert_json_value "${tmpdir}/resolve.json" '.packet.river.nodes | length >= 3' "resolve returns River nodes"

curl -fsSL -X POST "${BASE_URL}/resolve?session_id=${SESSION_ID}" \
  -H "Content-Type: text/plain" \
  --data 'Using the previous publication River, what should stay in context?' \
  -D "${tmpdir}/followup.headers" \
  -o "${tmpdir}/followup.json"
assert_json_value "${tmpdir}/followup.json" '.tripcode == "HUT-RIVER-001"' "second-turn follow-up reuses session TripCode"

echo
echo "== public repository smoke =="
http_get "${PUBLIC_REPO_URL}" "${tmpdir}/repo.html" "${tmpdir}/repo.headers"
rg -q "aitrailblazer-ai-agent-publishing" "${tmpdir}/repo.html" || fail "public GitHub repo page did not load expected name"
pass "public GitHub repo reachable"

http_get "${PUBLIC_RAW_BASE}/LICENSE" "${tmpdir}/repo-license.txt" "${tmpdir}/repo-license.headers"
rg -q "MIT License" "${tmpdir}/repo-license.txt" || fail "public repo raw LICENSE missing MIT License"
pass "public repo MIT license visible"

http_get "${PUBLIC_RAW_BASE}/${DOC_PATH}" "${tmpdir}/repo-doc.html" "${tmpdir}/repo-doc.headers"
rg -q "Consolidated Public Documentation" "${tmpdir}/repo-doc.html" || fail "public repo docs missing consolidated docs marker"
pass "public repo consolidated docs visible"

echo
echo "submission readiness passed"
