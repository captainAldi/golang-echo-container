# build
FROM golang:alpine as builder
LABEL maintainer="renaldiyulvianda@yahoo.com"

RUN apk update && \
  apk add --no-cache git ca-certificates tzdata \
  # && apk add --no-cache curl \
  && update-ca-certificates

# Move to working directory /build
WORKDIR /golang-echo-container

# Copy go.mod & go.sum, run go mod download
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Download dependency using go mod
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o go-checker .

# distribute
FROM alpine:3.13

WORKDIR /golang-echo-container

# Create appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \    
  --disabled-password \    
  --gecos "" \    
  --home "/nonexistent" \    
  --shell "/sbin/nologin" \    
  --no-create-home \    
  --uid "${UID}" \    
  "${USER}"

COPY --from=builder --chown=appuser:appuser /golang-echo-container/go-checker /golang-echo-container

USER appuser:appuser

STOPSIGNAL SIGINT

# Export necessary port
EXPOSE 8000

# Command to run when starting the container
ENTRYPOINT ["./go-checker"]