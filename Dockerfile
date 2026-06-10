FROM golang:1.24-bookworm AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o /out/aitrailblazer-ai-agent-publishing ./cmd/server

FROM mongodb/mongodb-community-server:7.0-ubuntu2204

USER root
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl gnupg \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && npm install -g mongodb-mcp-server@latest \
    && npm cache clean --force \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/aitrailblazer-ai-agent-publishing /aitrailblazer-ai-agent-publishing
COPY --from=build /src/index.html /index.html
COPY --from=build /src/README.html /README.html
COPY --from=build /src/CHANGELOG.html /CHANGELOG.html
COPY --from=build /src/DEVPOST_SUBMISSION.html /DEVPOST_SUBMISSION.html
COPY --from=build /src/EXECPLAN_FINAL_SUBMISSION_VIDEO.html /EXECPLAN_FINAL_SUBMISSION_VIDEO.html
COPY --from=build /src/LICENSE /LICENSE
COPY --from=build /src/img /img
COPY --from=build /src/docs /docs
COPY scripts/cloud-run-entrypoint.sh /cloud-run-entrypoint.sh
ENV PORT=8080
ENV MONGODB_DATABASE=aitrailblazer_demo
ENV MONGODB_DEPLOYMENT_KIND=embedded
ENV MONGODB_MCP_HOST=127.0.0.1
ENV MONGODB_MCP_PORT=3000
ENV MCP_SERVER_URL=http://127.0.0.1:3000/mcp
ENV MCP_METHOD=tools/call
ENV MCP_TOOL_NAME=find
ENV MCP_SESSION_ID=aitrailblazer-judge-proof
ENV MDB_MCP_CONNECTION_STRING=mongodb://127.0.0.1:27017/?directConnection=true
ENV MDB_MCP_READ_ONLY=true
ENV MDB_MCP_TELEMETRY=disabled
ENV MDB_MCP_EXTERNALLY_MANAGED_SESSIONS=true
ENV MDB_MCP_HTTP_RESPONSE_TYPE=json
EXPOSE 8080
ENTRYPOINT ["/cloud-run-entrypoint.sh"]
