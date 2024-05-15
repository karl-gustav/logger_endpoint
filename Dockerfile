FROM golang:1-alpine as builder

RUN apk update
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build.
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . ./

# CGO_ENABLED=0 == Don't depend on libc (bigger but more independent binary)
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM scratch
WORKDIR /app

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/main .

CMD ["./main"]
