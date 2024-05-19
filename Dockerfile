FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

WORKDIR /build

ADD go.mod .
COPY test_file.txt ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /build/task task.go

FROM alpine

WORKDIR /build
COPY --from=builder /build/task /build/task

CMD ["./task", "test_file.txt"]
