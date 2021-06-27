FROM golang:alpine AS builder

WORKDIR /src

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/service github.com/exelr/filotto/cmd/service

FROM alpine

RUN addgroup -S service && \
    adduser -S service -G service

WORKDIR /app
RUN chown -R service:service /app

USER service

COPY --chown=service:service --from=builder /app/service ./service

ENTRYPOINT ["./service"]
