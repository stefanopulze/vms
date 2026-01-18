# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -a \
    -installsuffix cgo \
    -o /app/dist/vms-core ./cmd/vms-core

FROM gcr.io/distroless/static-debian13:nonroot
ARG TARGETARCH

ENV TZ=Europe/Rome

WORKDIR /app
COPY --from=builder /app/dist/vms-core /app/vms-core

EXPOSE 8080

ENTRYPOINT ["/app/vms-core"]