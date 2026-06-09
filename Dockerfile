FROM golang:1.24-bookworm AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o /out/aitrailblazer-ai-agent-publishing ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/aitrailblazer-ai-agent-publishing /aitrailblazer-ai-agent-publishing
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/aitrailblazer-ai-agent-publishing"]
