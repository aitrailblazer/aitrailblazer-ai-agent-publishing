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

if rg -n "href=[\"']\\./docs/" index.html PROJECT_README.html CHANGELOG.html DEVPOST_SUBMISSION.html EXECPLAN_FINAL_SUBMISSION_VIDEO.html START_HERE.html VIDEO_SLIDES.html | rg -v "${DOC_PATH}" >/dev/null; then
  rg -n "href=[\"']\\./docs/" index.html PROJECT_README.html CHANGELOG.html DEVPOST_SUBMISSION.html EXECPLAN_FINAL_SUBMISSION_VIDEO.html START_HERE.html VIDEO_SLIDES.html | rg -v "${DOC_PATH}" >&2
  fail "public HTML links to stale docs files"
fi
pass "index/README/changelog/submission docs links target consolidated public docs"

echo
echo "== public secret and evidence scan =="
private_evidence_pattern="Co""dex|Open""AI|Chat""GPT|Clau""de|Anth""ropic|Cur""sor|Co""pilot|csrf""-token|Rapid_Agent_Hackathon_Rules_Source_2026_06_08""\\.raw\\.html|860""1137|AI""za|s""k-[A-Za-z0-9]|mong""odb\\+srv://[^ <)\"']+:[^ <)\"']+@"
if rg -n "${private_evidence_pattern}" \
  index.html PROJECT_README.html CHANGELOG.html DEVPOST_SUBMISSION.html EXECPLAN_FINAL_SUBMISSION_VIDEO.html START_HERE.html VIDEO_SLIDES.html Makefile Dockerfile cmd internal scripts docs .env.example >"${tmpdir}/scan.txt"; then
  cat "${tmpdir}/scan.txt" >&2
  fail "public scan found prohibited private evidence or obvious token pattern"
fi
pass "public scan found no prohibited private evidence terms or obvious token patterns"

echo
echo "== hosted HTTP smoke =="
http_get "${BASE_URL}/" "${tmpdir}/root.html" "${tmpdir}/root.headers"
rg -q "AITrailblazer AI Agent Publishing" "${tmpdir}/root.html" || fail "hosted root missing project title"
rg -q "Judge Proof Checklist" "${tmpdir}/root.html" || fail "hosted root missing judge proof checklist"
rg -q "Run Judge Demo" "${tmpdir}/root.html" || fail "hosted root missing Run Judge Demo button"
rg -q "Public docs" "${tmpdir}/root.html" || fail "hosted root missing consolidated docs link"
rg -q "Devpost Submission Pack" "${tmpdir}/root.html" || fail "hosted root missing Devpost submission pack link"
rg -q "Final Video ExecPlan" "${tmpdir}/root.html" || fail "hosted root missing Final Video ExecPlan link"
rg -q "Start Here" "${tmpdir}/root.html" || fail "hosted root missing Start Here link"
rg -q "Video Slides" "${tmpdir}/root.html" || fail "hosted root missing Video Slides link"
rg -q "A short code printed in an article" "${tmpdir}/root.html" || fail "hosted root missing plain-language TripCode definition"
rg -q "Resolve handle" "${tmpdir}/root.html" || fail "hosted root missing agent action trace"
rg -q "Mongo DB MCP Runtime Proof" "${tmpdir}/root.html" || fail "hosted root missing MongoDB MCP runtime proof panel"
pass "hosted root"

http_get "${BASE_URL}/health" "${tmpdir}/health.json" "${tmpdir}/health.headers"
assert_json_value "${tmpdir}/health.json" '.ok == true' "hosted health"

http_get "${BASE_URL}/${DOC_PATH}" "${tmpdir}/docs.html" "${tmpdir}/docs.headers"
rg -q "Consolidated Public Documentation" "${tmpdir}/docs.html" || fail "hosted docs missing consolidated docs marker"
rg -q "Embedded XML Contract" "${tmpdir}/docs.html" || fail "hosted docs missing embedded XML marker"
rg -q "Judge Quick Proof" "${tmpdir}/docs.html" || fail "hosted docs missing judge quick proof marker"
pass "hosted consolidated docs"

http_get "${BASE_URL}/LICENSE" "${tmpdir}/license.txt" "${tmpdir}/license.headers"
rg -q "MIT License" "${tmpdir}/license.txt" || fail "hosted license missing MIT License"
pass "hosted license"

http_get "${BASE_URL}/DEVPOST_SUBMISSION.html" "${tmpdir}/devpost.html" "${tmpdir}/devpost.headers"
rg -q "Copy-ready Devpost submission pack" "${tmpdir}/devpost.html" || fail "hosted Devpost submission pack missing title marker"
rg -q "Demo Video URL: Add final under-3-minute demo video link before submission" "${tmpdir}/devpost.html" || fail "hosted Devpost submission pack missing video placeholder"
pass "hosted Devpost submission pack"

http_get "${BASE_URL}/EXECPLAN_FINAL_SUBMISSION_VIDEO.html" "${tmpdir}/video-execplan.html" "${tmpdir}/video-execplan.headers"
rg -q "Final Submission Video ExecPlan" "${tmpdir}/video-execplan.html" || fail "hosted Final Video ExecPlan missing title"
rg -q "make video-url-check DEMO_VIDEO_URL" "${tmpdir}/video-execplan.html" || fail "hosted Final Video ExecPlan missing video URL check"
pass "hosted Final Video ExecPlan"

http_get "${BASE_URL}/START_HERE.html" "${tmpdir}/start-here.html" "${tmpdir}/start-here.headers"
rg -q "AITrailblazer Start Here" "${tmpdir}/start-here.html" || fail "hosted Start Here missing title"
rg -q "Replay Steps" "${tmpdir}/start-here.html" || fail "hosted Start Here missing replay steps"
rg -q "HUT-RIVER-001" "${tmpdir}/start-here.html" || fail "hosted Start Here missing TripCode"
pass "hosted Start Here"

http_get "${BASE_URL}/VIDEO_SLIDES.html" "${tmpdir}/video-slides.html" "${tmpdir}/video-slides.headers"
rg -q "AITrailblazer Video Slides" "${tmpdir}/video-slides.html" || fail "hosted Video Slides missing title"
rg -q "Under-three-minute public recording storyboard" "${tmpdir}/video-slides.html" || fail "hosted Video Slides missing storyboard marker"
pass "hosted Video Slides"

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

http_get "${PUBLIC_RAW_BASE}/PROJECT_README.html" "${tmpdir}/repo-readme.html" "${tmpdir}/repo-readme.headers"
rg -q "Judge Quickstart" "${tmpdir}/repo-readme.html" || fail "public repo README missing Judge Quickstart"
rg -q "Public Safety Boundary" "${tmpdir}/repo-readme.html" || fail "public repo README missing Public Safety Boundary"
rg -q "No investment recommendations" "${tmpdir}/repo-readme.html" || fail "public repo README missing investment boundary"
pass "public repo README judge quickstart visible"

http_get "${PUBLIC_RAW_BASE}/DEVPOST_SUBMISSION.html" "${tmpdir}/repo-devpost.html" "${tmpdir}/repo-devpost.headers"
rg -q "AITrailblazer AI Agent Publishing Devpost Submission Pack" "${tmpdir}/repo-devpost.html" || fail "public repo Devpost submission pack missing title"
rg -q "does not claim to be a live SEC/XBRL market-data product" "${tmpdir}/repo-devpost.html" || fail "public repo Devpost submission pack missing scope boundary"
pass "public repo Devpost submission pack visible"

http_get "${PUBLIC_RAW_BASE}/EXECPLAN_FINAL_SUBMISSION_VIDEO.html" "${tmpdir}/repo-video-execplan.html" "${tmpdir}/repo-video-execplan.headers"
rg -q "AITrailblazer Final Submission Video ExecPlan" "${tmpdir}/repo-video-execplan.html" || fail "public repo Final Video ExecPlan missing title"
pass "public repo Final Video ExecPlan visible"

http_get "${PUBLIC_RAW_BASE}/START_HERE.html" "${tmpdir}/repo-start-here.html" "${tmpdir}/repo-start-here.headers"
rg -q "AITrailblazer Start Here" "${tmpdir}/repo-start-here.html" || fail "public repo Start Here missing title"
pass "public repo Start Here visible"

http_get "${PUBLIC_RAW_BASE}/VIDEO_SLIDES.html" "${tmpdir}/repo-video-slides.html" "${tmpdir}/repo-video-slides.headers"
rg -q "AITrailblazer Video Slides" "${tmpdir}/repo-video-slides.html" || fail "public repo Video Slides missing title"
pass "public repo Video Slides visible"

http_get "${PUBLIC_RAW_BASE}/${DOC_PATH}" "${tmpdir}/repo-doc.html" "${tmpdir}/repo-doc.headers"
rg -q "Consolidated Public Documentation" "${tmpdir}/repo-doc.html" || fail "public repo docs missing consolidated docs marker"
rg -q "Mongo DB MCP Runtime Proof" "${tmpdir}/repo-doc.html" || fail "public repo docs missing Mongo DB MCP Runtime Proof"
pass "public repo consolidated docs visible"

echo
echo "submission readiness passed"
