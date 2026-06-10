#!/usr/bin/env bash
set -euo pipefail

fail() {
  printf 'FAIL - %s\n' "$1" >&2
  exit 1
}

pass() {
  printf 'ok - %s\n' "$1"
}

url="${DEMO_VIDEO_URL:-${1:-}}"
[[ -n "${url}" ]] || fail "set DEMO_VIDEO_URL to the final public YouTube or Vimeo demo URL"

case "${url}" in
  https://www.youtube.com/*|https://youtu.be/*|https://youtube.com/*|https://vimeo.com/*|https://www.vimeo.com/*)
    ;;
  *)
    fail "video URL must be a public YouTube or Vimeo HTTPS URL"
    ;;
esac

headers="$(mktemp)"
body="$(mktemp)"
cleanup() {
  rm -f "${headers}" "${body}"
}
trap cleanup EXIT

curl -fsSL -D "${headers}" -o "${body}" "${url}" || fail "video URL did not return a successful response"

content_type="$(tr -d '\r' <"${headers}" | awk 'BEGIN{IGNORECASE=1} /^content-type:/ {print $2; exit}')"
case "${content_type}" in
  text/html*|application/xhtml+xml*)
    ;;
  *)
    fail "video URL returned unexpected content type: ${content_type:-missing}"
    ;;
esac

if [[ "${url}" == *"youtube"* || "${url}" == *"youtu.be"* ]]; then
  rg -qi "youtube|ytInitialPlayerResponse|watch" "${body}" || fail "YouTube page marker not found"
else
  rg -qi "vimeo|player" "${body}" || fail "Vimeo page marker not found"
fi

pass "public demo video URL is reachable"
