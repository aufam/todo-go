FROM golang:1.23-alpine3.20 AS builder

WORKDIR /root/todo

COPY go.mod go.sum ./
RUN go mod download

COPY app app/
COPY api api/
COPY core core/
COPY db db/
COPY models models/
COPY tests tests/

RUN go test ./tests && \
    go build -o todo app/main.go

FROM alpine:3.20 AS runner

WORKDIR /root
COPY --from=builder /root/todo/todo .

EXPOSE 8000
CMD ["./todo"]
