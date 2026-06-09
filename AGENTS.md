# AGENTS.md

This file governs agent work in the `aitrailblazer-ai-agent-publishing` workspace.

## Project Boundary

AITrailblazer AI Agent Publishing is the Google Cloud Rapid Agent Hackathon workspace for packaging a publication-to-agent-memory product.

The core idea starts from DeltaSignal: a Substack publication with more than 300 macro and financial articles introduced TripCodes and Rivers so one article can become a stable, agent-resolvable research object. This workspace generalizes that pattern for other publishers, analysts, creators, and research teams that want their own archives to become agent-readable memory.

This repository is not the secret-bearing runtime. It owns public-facing submission materials, static visual specs, demo copy, judge checklists, and links to runnable agent surfaces.

## Output Rules

- Durable long-form artifacts must be self-contained StrategiX Visual Spec `.html` files with an embedded XML contract.
- Do not create parallel `.md` specifications, roadmaps, reports, or prompt packets.
- Keep Devpost, Google Cloud, Gemini, Agent Builder, Cloud Run, MongoDB, partner, and evidence claims grounded in current source material.
- Never commit API keys, Devpost session data, Google Cloud secrets, billing identifiers that are not intentionally public, or private account screenshots.

## Local Commands

Use these from this repository root:

```bash
GOWORK=off go test ./...
make check
make run
curl http://localhost:8080/health
curl -s http://localhost:8080/v1/archive-brief \
  -H 'Content-Type: application/json' \
  -d '{"publication_url":"https://deltasignal.substack.com/"}'
curl -s http://localhost:8080/v1/tripcode \
  -H 'Content-Type: application/json' \
  -d '{"tripcode":"AIT-PUB-HUT-8-RERATE","session_id":"demo"}'
make validate
make preview
```

The public static entrypoint is `index.html`.

The public demo endpoints are `POST /v1/archive-brief`, `POST /v1/tripcode`, judge-friendly `GET/POST /resolve`, `GET /v1/usage`, and `GET /health`.

## Publishing Workflow

1. Keep `index.html` judge-readable and GitHub Pages-compatible.
2. Put long-form packets in `docs/*.html`.
3. Keep generated or downloaded proof artifacts in `artifacts/` if needed; do not commit private screenshots unless explicitly approved.
4. Use the Rapid Agent resources page as the canonical public resource map for Google Cloud agent setup, partner integrations, recordings, FAQ, forum, and Discord links.
5. Treat MongoDB as the closest partner fit only when discussing archive memory, vector search, semantic retrieval, MCP access, and operational persistence.
6. Use Atlas Free for no-cost MongoDB development and smoke tests. If paid Atlas usage is needed for contest work, subscribe through Google Cloud Marketplace so usage is unified under Google Cloud billing and can draw down hackathon credits.
7. Treat the MongoDB Gemini CLI extension as a free permitted helper, subject to current Gemini CLI requirements and Google account setup.
8. Keep runtime claims separate from publishing claims: this workspace can point to deployed services, but it should not own backend credentials or call protected APIs.

## Competition Boundary

The publishing surface may describe:

- Google Cloud Rapid Agent Hackathon positioning.
- Gemini Enterprise Agent Platform / Vertex AI / Agent Builder setup.
- Agent Starter Pack, Python SDK, Agent Runtime, Secret Manager, and Cloud Run paths.
- Partner resource lanes such as Arize, Elastic, Fivetran, GitLab, MongoDB, and Dynatrace.
- The DeltaSignal proof case: a 300+ article Substack archive, HUT TripCode, and River memory.
- The generalized submission narrative: publishers can convert their own articles into TripCodes, Rivers, vector/semantic memory, and agent workflows.
- The AITrailblazer submission narrative, demo flow, architecture, and judge checklist.

It must not claim a deployed runtime, score, sponsor approval, or partner integration is live unless the linked project or current evidence confirms it.

## Final Submission Gate

Before Devpost submission, verify:

- Hosted project URL opens from a fresh browser and points to the running app, not the repository.
- Public repository opens while signed out and includes an OSI-approved root `LICENSE` file.
- README has clear setup and run instructions.
- Demo video is public, on YouTube or Vimeo, and under three minutes.
- MongoDB track is selected and all team members are listed.
- Gemini, Google Cloud Agent Builder, and MongoDB MCP server are invoked at runtime, not only named in docs.
- Final contest-facing implementation and evidence contain no competing AI or cloud services.

## AI Tool Compliance

The contest rules and Devpost FAQ create a strict AI-tool boundary for both the final project and the development workflow.

- Permitted AI tools: Google Cloud AI tools such as Gemini models, Google Cloud Agent Builder, and Google's AntiGravity suite; selected track partner built-in AI features.
- Not permitted for final project dependencies or public runtime evidence: unrelated non-Google AI services or competing cloud services listed as prohibited by the official rules or FAQ.
- Before preparing the public submission repository, audit source, docs, lockfiles, scripts, screenshots, environment files, and commit metadata for competitor AI tool dependencies or development traces.
- Keep private local drafts and operator traces out of the public submission repository.

## Repository Contents

- `index.html`: public static landing/spec page for GitHub Pages or local preview.
- `CHANGELOG.html`: public changelog.
- `README.html`: local workspace notes.
- `docs/Rapid_Agent_Hackathon_Resource_Map_2026_06_08.html`: Rapid Agent resource map plus MongoDB fit.
- `assets/agent-publishing-evidence-pipeline.png`: visual project signal borrowed from the DeltaSignal proof case.
- `cmd/server`: HTTP service for judge demo endpoints.
- `internal/agent`: coordinator, TripCode resolver, memory store, and cost tracker.
- `scripts/judge-demo.sh`: curl flow for health, archive brief, TripCode resolution, and session memory.
- `scripts/report-spend.sh`: usage ledger check.
- `Dockerfile`: Cloud Run-ready container image.
- `cloudbuild.yaml`: Cloud Build and Cloud Run deployment template.
- `Makefile`: local test/build/demo/deploy shortcuts.
