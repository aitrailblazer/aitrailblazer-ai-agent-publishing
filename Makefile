APP_NAME ?= aitrailblazer-ai-agent-publishing
SERVICE ?= aitrailblazer-ai-agent-publishing
PROJECT ?= aitrailblazer-ai-publishing
REGION ?= us-central1
IMAGE ?= $(REGION)-docker.pkg.dev/$(PROJECT)/$(APP_NAME)/app:local
PORT ?= 8080
ATLAS_SECRET ?= publishing-mongodb-atlas-uri

-include .env
export

.PHONY: help test coverage vet check run demo spend build docker-build deploy deploy-source deploy-atlas preview validate submission-readiness mcp-check atlas-check mcp-run clean print-config

help:
	@echo "Available commands:"
	@echo "  make test          - Run Go tests"
	@echo "  make coverage      - Run Go tests with coverage report"
	@echo "  make vet           - Run go vet"
	@echo "  make check         - Run tests, coverage, vet, and static validation"
	@echo "  make run           - Run the Cloud Run HTTP service locally"
	@echo "  make demo          - Run judge-friendly curl demo flow"
	@echo "  make spend         - Report app usage ledger"
	@echo "  make build         - Build local server binary"
	@echo "  make docker-build  - Build local Docker image"
	@echo "  make deploy        - Deploy through cloudbuild.yaml"
	@echo "  make deploy-source - Deploy source directly to Cloud Run"
	@echo "  make deploy-atlas  - Build and deploy Cloud Run with Atlas URI from Secret Manager"
	@echo "  make preview   - Serve the static publishing surface locally"
	@echo "  make validate  - Run lightweight static validation"
	@echo "  make submission-readiness - Run final local, hosted, repo, link, and scan gates"
	@echo "  make mcp-check - Verify official MongoDB MCP prerequisites"
	@echo "  make atlas-check - Verify Atlas-mode MongoDB MCP prerequisites"
	@echo "  make mcp-run   - Run official MongoDB MCP server over HTTP"
	@echo "  make clean     - Remove local generated files"

test:
	GOWORK=off go test ./...

coverage:
	GOWORK=off go test ./... -coverprofile=coverage.out -covermode=count
	GOWORK=off go tool cover -func=coverage.out
	@GOWORK=off go tool cover -func=coverage.out | awk '/^total:/ { if ($$3 != "100.0%") { print "coverage gate failed: " $$3; exit 1 } }'

vet:
	GOWORK=off go vet ./...

check: test coverage vet validate

run:
	PORT=$(PORT) GOWORK=off go run ./cmd/server

demo:
	scripts/judge-demo.sh

spend:
	scripts/report-spend.sh

build:
	mkdir -p bin
	CGO_ENABLED=0 GOWORK=off go build -trimpath -ldflags="-w -s" -o bin/$(APP_NAME) ./cmd/server

docker-build:
	docker build -t $(APP_NAME):latest .

deploy:
	gcloud builds submit \
		--project=$(PROJECT) \
		--config=cloudbuild.yaml \
		--substitutions=_REGION=$(REGION),_IMAGE=$(IMAGE)

deploy-source:
	gcloud run deploy $(SERVICE) \
		--source . \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed \
		--allow-unauthenticated \
		--min-instances=0 \
		--max-instances=3 \
		--memory=1Gi \
		--cpu=1 \
		--concurrency=80 \
		--set-env-vars=GOOGLE_CLOUD_PROJECT=$(PROJECT),GOOGLE_CLOUD_LOCATION=us-central1,GOOGLE_GENAI_USE_VERTEXAI=true,PUBLISHING_USE_GEMINI=true,GEMINI_MODEL=gemini-2.5-flash,AGENT_BUILDER_MODE=discoveryengine-search,AGENT_BUILDER_ENDPOINT=https://discoveryengine.googleapis.com/v1/projects/$(PROJECT)/locations/global/collections/default_collection/engines/aitrailblazer-pub-agent/servingConfigs/default_search:search,AGENT_BUILDER_AGENT_ID=aitrailblazer-pub-agent,MONGODB_DATABASE=aitrailblazer_demo,MONGODB_DEPLOYMENT_KIND=embedded,MCP_SERVER_URL=http://127.0.0.1:3000/mcp,MCP_METHOD=tools/call,MCP_TOOL_NAME=find,MCP_SESSION_ID=aitrailblazer-judge-proof,MDB_MCP_CONNECTION_STRING=mongodb://127.0.0.1:27017/?directConnection=true,MDB_MCP_READ_ONLY=true,MDB_MCP_TELEMETRY=disabled,MDB_MCP_EXTERNALLY_MANAGED_SESSIONS=true,MDB_MCP_HTTP_RESPONSE_TYPE=json,PUBLISHING_COST_TRACKING=true,PUBLISHING_GOOGLE_CREDIT_BUDGET_USD=500,PUBLISHING_ESTIMATED_ARCHIVE_BRIEF_COST_USD=0.01,PUBLISHING_ESTIMATED_TRIPCODE_COST_USD=0.02,PUBLISHING_ESTIMATED_JUDGE_DEMO_COST_USD=0.02,PUBLISHING_ESTIMATED_SESSION_MEMORY_COST_USD=0.00,PUBLISHING_COST_SOURCE=local-estimate

deploy-atlas:
	gcloud builds submit \
		--project=$(PROJECT) \
		--tag=$(IMAGE) .
	gcloud run deploy $(SERVICE) \
		--image=$(IMAGE) \
		--project=$(PROJECT) \
		--region=$(REGION) \
		--platform=managed \
		--allow-unauthenticated \
		--min-instances=0 \
		--max-instances=3 \
		--memory=1Gi \
		--cpu=1 \
		--concurrency=80 \
		--set-secrets=MDB_MCP_CONNECTION_STRING=$(ATLAS_SECRET):latest \
		--set-env-vars=GOOGLE_CLOUD_PROJECT=$(PROJECT),GOOGLE_CLOUD_LOCATION=us-central1,GOOGLE_GENAI_USE_VERTEXAI=true,PUBLISHING_USE_GEMINI=true,GEMINI_MODEL=gemini-2.5-flash,AGENT_BUILDER_MODE=discoveryengine-search,AGENT_BUILDER_ENDPOINT=https://discoveryengine.googleapis.com/v1/projects/$(PROJECT)/locations/global/collections/default_collection/engines/aitrailblazer-pub-agent/servingConfigs/default_search:search,AGENT_BUILDER_AGENT_ID=aitrailblazer-pub-agent,MONGODB_DATABASE=aitrailblazer_demo,MONGODB_DEPLOYMENT_KIND=atlas,MONGODB_SEED_DEMO_DATA=true,MCP_SERVER_URL=http://127.0.0.1:3000/mcp,MCP_METHOD=tools/call,MCP_TOOL_NAME=find,MCP_SESSION_ID=aitrailblazer-judge-proof,MDB_MCP_READ_ONLY=true,MDB_MCP_TELEMETRY=disabled,MDB_MCP_EXTERNALLY_MANAGED_SESSIONS=true,MDB_MCP_HTTP_RESPONSE_TYPE=json,PUBLISHING_COST_TRACKING=true,PUBLISHING_GOOGLE_CREDIT_BUDGET_USD=500,PUBLISHING_ESTIMATED_ARCHIVE_BRIEF_COST_USD=0.01,PUBLISHING_ESTIMATED_TRIPCODE_COST_USD=0.02,PUBLISHING_ESTIMATED_JUDGE_DEMO_COST_USD=0.02,PUBLISHING_ESTIMATED_SESSION_MEMORY_COST_USD=0.00,PUBLISHING_COST_SOURCE=local-estimate

preview:
	python3 -m http.server 8097

validate:
	@test -f index.html
	@test -f img/AITrailblazerAI.png
	@test -f README.html
	@test -f DEVPOST_SUBMISSION.html
	@test -f LICENSE
	@test -f AGENTS.md
	@test -f go.mod
	@test -f cmd/server/main.go
	@test -f internal/agent/types.go
	@test -f docs/AITrailblazer_AI_Agent_Publishing_Public_Docs_2026_06_10.html
	@test "$$(find docs -maxdepth 1 -type f -name '*.html' | wc -l | tr -d ' ')" = "1"
	@grep -q "StrategiXVisualSpec" index.html
	@grep -q "StrategiXVisualSpec" DEVPOST_SUBMISSION.html
	@grep -q "StrategiXVisualSpec" docs/AITrailblazer_AI_Agent_Publishing_Public_Docs_2026_06_10.html
	@grep -q "Google Cloud Rapid Agent Hackathon" index.html
	@grep -q "rapid-agent.devpost.com/resources" index.html
	@grep -q "AITrailblazer AI Agent Publishing" index.html
	@echo "Static publishing surface validated."

submission-readiness:
	scripts/submission-readiness.sh

mcp-check:
	scripts/mongodb-mcp-check.sh

atlas-check:
	scripts/mongodb-atlas-check.sh

mcp-run:
	scripts/mongodb-mcp-run-http.sh

print-config:
	@echo "APP_NAME=$(APP_NAME)"
	@echo "SERVICE=$(SERVICE)"
	@echo "PROJECT=$(PROJECT)"
	@echo "REGION=$(REGION)"
	@echo "IMAGE=$(IMAGE)"
	@echo "PORT=$(PORT)"

clean:
	rm -rf bin/ tmp/ coverage.out
	docker rmi $(APP_NAME):latest 2>/dev/null || true
