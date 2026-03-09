# syntax=docker/dockerfile:1
### Frontend
FROM oven/bun as frontend

WORKDIR /app
COPY vms-ui/package*.json ./
COPY vms-ui/bun*.lock ./

RUN bun install --frozen-lockfile

COPY ./vms-ui/ .
RUN bun run build-only

### Backend
FROM --platform=$BUILDPLATFORM golang:1.26 AS backend
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY ./vms-core/go.mod ./vms-core/go.sum* ./
RUN go mod tidy && go mod download

COPY ./vms-core/ .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -a \
    -installsuffix cgo \
    -o /app/dist/vms-core ./cmd/vms-core

### Final
FROM gcr.io/distroless/static-debian13:nonroot
ARG TARGETARCH

ENV TZ=Europe/Rome

WORKDIR /app
COPY --from=backend /app/dist/vms-core /app/vms-core
COPY --from=frontend /app/dist/ /app/web

EXPOSE 8080

ENTRYPOINT ["/app/vms-core"]