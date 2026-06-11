# AITrailblazer AI Agent Publishing

Google Cloud Rapid Agent Hackathon submission for turning a publication archive into agent-readable memory.

The project demonstrates a browser and API workflow where a published research archive becomes structured agent context: article objects, TripCodes, River continuity, MongoDB-backed state, Gemini-ready synthesis, and reusable reader packets.

## Live Demo

- Hosted app: https://aitrailblazer-ai-agent-publishing-rmycwek6ba-uc.a.run.app/
- Browser demo: https://aitrailblazer-ai-agent-publishing-rmycwek6ba-uc.a.run.app/demo/
- Manual demo controls: https://aitrailblazer-ai-agent-publishing-rmycwek6ba-uc.a.run.app/demo/run

## What It Proves

- A Cloud Run application can expose a repeatable publishing-agent workflow.
- A TripCode can resolve a publication object into reusable River memory.
- MongoDB/MCP-shaped context can support article, claim, River, session, and agent-run records.
- Gemini-ready synthesis can convert raw JSON proof into a reader-facing packet.
- Usage visibility keeps the demo cost-aware.

## Run Locally

```bash
GOWORK=off go test ./...
PORT=8080 GOWORK=off go run ./cmd/server
```

Then open:

```text
http://127.0.0.1:8080/
http://127.0.0.1:8080/demo/
http://127.0.0.1:8080/demo/run
```

## Useful Commands

```bash
make test
make validate
make run
make deploy
npm ci
npm run check:playwright
```

## Main Routes

- `GET /health`
- `GET /v1/judge-demo`
- `GET /v1/usage`
- `POST /v1/archive-brief`
- `POST /v1/tripcode`
- `GET|POST /resolve`
- `GET /demo/`
- `GET /demo/run`

## Public Docs

The durable public project documentation is maintained as self-contained HTML visual specs:

- `index.html`
- `README.html`
- `START_HERE.html`
- `VIDEO_SLIDES.html`
- `docs/AITrailblazer_AI_Agent_Publishing_Public_Docs_2026_06_10.html`

No API keys, private screenshots, or local operator artifacts are required to run the public demo paths.
