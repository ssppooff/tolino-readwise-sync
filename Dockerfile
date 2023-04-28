# syntax=docker/dockerfile:1

## Build source
FROM golang:1.20.3-alpine AS build-stage
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY *.go ./
COPY readwise/*.go ./readwise/
COPY tolino/*.go ./tolino/
COPY utils/*.go ./utils/

# Build and rename binary to tolino-sync
RUN CGO_ENABLED=0 GOOS=linux go build -o /tolino-sync

## Run tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

## Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /tolino-sync /tolino-sync

LABEL org.opencontainers.image.source="https://github.com/ssppooff/tolino-readwise-sync"

USER nonroot:nonroot

ENTRYPOINT ["/tolino-sync"]
CMD ["-t", "/files/token", "-n", "/files/notes.txt"]
