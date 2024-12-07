FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod install \
  && go mod tidy \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main main.go

FROM alpine:latest

RUN addgroup -g 1001 -S appuser && adduser -u 1001 -S appuser -G appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser ./app/main /app/main

USER appuser

CMD [ "main" ]