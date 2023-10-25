# Step 1: Modules caching
FROM golang:1.21-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/user

# Step 3: Final
FROM alpine:latest

RUN apk add --no-cache tzdata

# GOPATH for scratch images is /
COPY --from=builder /bin/app /app
RUN mkdir config
CMD ["/app"]